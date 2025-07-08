package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"atempo/internal/utils"
)

// MCPClient handles communication with MCP servers
type MCPClient struct {
	ServerPath  string
	ProjectPath string
	Process     *exec.Cmd
	StdinPipe   io.WriteCloser
	StdoutPipe  io.ReadCloser
	RequestID   int
}

// MCPRequest represents a JSON-RPC request to an MCP server
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// MCPResponse represents a JSON-RPC response from an MCP server
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error in MCP communication
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TestParams represents parameters for test tool calls
type TestParams struct {
	Path    string   `json:"path,omitempty"`
	Options []string `json:"options,omitempty"`
}

// NewMCPClient creates a new MCP client for a project
func NewMCPClient(projectPath string) (*MCPClient, error) {
	// Look for MCP server in project
	mcpServerPath := filepath.Join(projectPath, "ai", "mcp-server", "index.js")
	if !utils.FileExists(mcpServerPath) {
		return nil, fmt.Errorf("no MCP server found at %s", mcpServerPath)
	}

	return &MCPClient{
		ServerPath:  mcpServerPath,
		ProjectPath: projectPath,
		RequestID:   1,
	}, nil
}

// Start launches the MCP server process
func (c *MCPClient) Start() error {
	// Start the MCP server process
	c.Process = exec.Command("node", c.ServerPath)
	c.Process.Dir = c.ProjectPath

	// Set up pipes for communication
	var err error
	c.StdinPipe, err = c.Process.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	c.StdoutPipe, err = c.Process.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the process
	if err := c.Process.Start(); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}

// Close terminates the MCP server and cleans up resources
func (c *MCPClient) Close() error {
	if c.StdinPipe != nil {
		c.StdinPipe.Close()
	}
	if c.StdoutPipe != nil {
		c.StdoutPipe.Close()
	}
	if c.Process != nil {
		return c.Process.Process.Kill()
	}
	return nil
}

// GetTestCommand queries the MCP server for the appropriate test command
func (c *MCPClient) GetTestCommand(testSuite string) (string, error) {
	// First, list available tools to see what test tools are available
	tools, err := c.ListTools()
	if err != nil {
		return "", fmt.Errorf("failed to list MCP tools: %w", err)
	}

	// Look for test-related tools
	var testTool string
	for _, tool := range tools {
		if name, ok := tool["name"].(string); ok {
			if strings.Contains(name, "test") {
				testTool = name
				break
			}
		}
	}

	if testTool == "" {
		return "", fmt.Errorf("no test tool found in MCP server")
	}

	// Prepare test parameters
	params := TestParams{}
	if testSuite != "" {
		params.Path = testSuite
	}

	// Call the test tool to get the command (we'll simulate this for now)
	// In a real implementation, we might have a "get_test_command" tool
	switch testTool {
	case "laravel_test":
		if testSuite != "" {
			return fmt.Sprintf("php artisan test --filter=%s", testSuite), nil
		}
		return "php artisan test", nil
	case "django_test":
		if testSuite != "" {
			return fmt.Sprintf("python manage.py test %s", testSuite), nil
		}
		return "python manage.py test", nil
	default:
		return "", fmt.Errorf("unsupported test tool: %s", testTool)
	}
}

// RunTests executes tests via the MCP server
func (c *MCPClient) RunTests(testSuite string) error {
	// List available tools first
	tools, err := c.ListTools()
	if err != nil {
		return fmt.Errorf("failed to list MCP tools: %w", err)
	}

	// Find test tool
	var testTool string
	for _, tool := range tools {
		if name, ok := tool["name"].(string); ok {
			if strings.Contains(name, "test") {
				testTool = name
				break
			}
		}
	}

	if testTool == "" {
		return fmt.Errorf("no test tool found in MCP server")
	}

	// Prepare test parameters
	params := TestParams{}
	if testSuite != "" {
		params.Path = testSuite
	}

	// Execute the test tool
	response, err := c.CallTool(testTool, params)
	if err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	// Display the test results
	if response.Error != nil {
		return fmt.Errorf("test tool error: %s", response.Error.Message)
	}

	// Extract and display content from response
	if result, ok := response.Result.(map[string]interface{}); ok {
		if content, ok := result["content"].([]interface{}); ok {
			for _, item := range content {
				if contentItem, ok := item.(map[string]interface{}); ok {
					if text, ok := contentItem["text"].(string); ok {
						fmt.Print(text)
					}
				}
			}
		}
	}

	return nil
}

// ListTools requests the list of available tools from the MCP server
func (c *MCPClient) ListTools() ([]map[string]interface{}, error) {
	request := MCPRequest{
		JSONRPC: "2.0",
		ID:      c.RequestID,
		Method:  "tools/list",
		Params:  map[string]interface{}{},
	}
	c.RequestID++

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("MCP error: %s", response.Error.Message)
	}

	// Extract tools from response
	if result, ok := response.Result.(map[string]interface{}); ok {
		if tools, ok := result["tools"].([]interface{}); ok {
			var toolList []map[string]interface{}
			for _, tool := range tools {
				if toolMap, ok := tool.(map[string]interface{}); ok {
					toolList = append(toolList, toolMap)
				}
			}
			return toolList, nil
		}
	}

	return nil, fmt.Errorf("unexpected response format")
}

// CallTool executes a specific tool on the MCP server
func (c *MCPClient) CallTool(toolName string, args interface{}) (*MCPResponse, error) {
	request := MCPRequest{
		JSONRPC: "2.0",
		ID:      c.RequestID,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		},
	}
	c.RequestID++

	return c.sendRequest(request)
}

// sendRequest sends a JSON-RPC request and waits for a response
func (c *MCPClient) sendRequest(request MCPRequest) (*MCPResponse, error) {
	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Add newline for line-delimited JSON
	requestJSON = append(requestJSON, '\n')

	// Send request
	if _, err := c.StdinPipe.Write(requestJSON); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	scanner := bufio.NewScanner(c.StdoutPipe)
	if scanner.Scan() {
		var response MCPResponse
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return &response, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return nil, fmt.Errorf("no response received")
}
