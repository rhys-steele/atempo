const server = Bun.serve({
  port: 3000,

  fetch(req) {
    const url = new URL(req.url);

    if (url.pathname === "/") {
      return Response.json({ message: "Hello from {{project}}" });
    }

    if (url.pathname === "/health") {
      return Response.json({ status: "ok" });
    }

    return new Response("Not Found", { status: 404 });
  },
});

console.log(`Listening on http://localhost:${server.port}`);
