# {{project}} - UI/UX Guidelines

## Design Principles

### 1. Modern, Clean Interface
- **Minimalist Design**: Clean layouts with plenty of whitespace
- **Consistent Typography**: Clear hierarchy with readable fonts
- **Intuitive Navigation**: Logical flow and clear call-to-actions
- **Responsive Design**: Mobile-first approach with progressive enhancement

### 2. Professional Aesthetic
- **Subtle Animations**: Smooth transitions that enhance user experience
- **Color Harmony**: Consistent color palette with accessible contrast ratios
- **Visual Hierarchy**: Clear distinction between primary and secondary elements
- **Brand Consistency**: Unified visual language across all interfaces

### 3. User-Centered Design
- **Accessibility First**: WCAG 2.1 AA compliance
- **Performance Focused**: Fast loading times and smooth interactions
- **Error Prevention**: Clear validation and helpful error messages
- **Feedback Loops**: Immediate response to user actions

## Frontend Architecture

### 1. Component Structure

#### Blade Components
```php
// resources/views/components/button.blade.php
@props([
    'type' => 'button',
    'variant' => 'primary',
    'size' => 'medium',
    'disabled' => false
])

<button 
    type="{{ $type }}"
    {{ $attributes->merge(['class' => "btn btn-{$variant} btn-{$size}"]) }}
    @if($disabled) disabled @endif
>
    {{ $slot }}
</button>
```

#### Usage Example
```blade
<x-button variant="primary" size="large" class="w-full">
    Create Account
</x-button>

<x-button variant="secondary" type="submit" :disabled="$form->processing">
    Save Changes
</x-button>
```

### 2. Layout System

#### Base Layout
```blade
<!-- resources/views/layouts/app.blade.php -->
<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="csrf-token" content="{{ csrf_token() }}">
    
    <title>{{ config('app.name', 'Laravel') }}</title>
    
    <!-- Styles -->
    @vite(['resources/css/app.css', 'resources/js/app.js'])
</head>
<body class="font-sans antialiased bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
        @include('partials.navigation')
    </nav>
    
    <!-- Main Content -->
    <main class="py-4">
        {{ $slot }}
    </main>
    
    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-auto">
        @include('partials.footer')
    </footer>
</body>
</html>
```

### 3. Form Design Patterns

#### Form Components
```blade
<!-- resources/views/components/form/input.blade.php -->
@props([
    'label' => null,
    'name' => null,
    'type' => 'text',
    'required' => false,
    'error' => null
])

<div class="mb-4">
    @if($label)
        <label for="{{ $name }}" class="block text-sm font-medium text-gray-700 mb-2">
            {{ $label }}
            @if($required)
                <span class="text-red-500">*</span>
            @endif
        </label>
    @endif
    
    <input
        type="{{ $type }}"
        name="{{ $name }}"
        id="{{ $name }}"
        {{ $attributes->merge(['class' => 'w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent']) }}
        @if($required) required @endif
    >
    
    @if($error)
        <p class="mt-1 text-sm text-red-600">{{ $error }}</p>
    @endif
</div>
```

## Visual Design System

### 1. Color Palette

#### Primary Colors
```css
:root {
    /* Primary Brand Colors */
    --color-primary-50: #eff6ff;
    --color-primary-100: #dbeafe;
    --color-primary-500: #3b82f6;
    --color-primary-600: #2563eb;
    --color-primary-700: #1d4ed8;
    
    /* Semantic Colors */
    --color-success: #10b981;
    --color-warning: #f59e0b;
    --color-error: #ef4444;
    --color-info: #06b6d4;
    
    /* Neutral Colors */
    --color-gray-50: #f9fafb;
    --color-gray-100: #f3f4f6;
    --color-gray-500: #6b7280;
    --color-gray-700: #374151;
    --color-gray-900: #111827;
}
```

#### CSS Implementation
```css
/* Button Variants */
.btn-primary {
    @apply bg-blue-600 hover:bg-blue-700 text-white;
}

.btn-secondary {
    @apply bg-gray-200 hover:bg-gray-300 text-gray-800;
}

.btn-success {
    @apply bg-green-600 hover:bg-green-700 text-white;
}

.btn-danger {
    @apply bg-red-600 hover:bg-red-700 text-white;
}
```

### 2. Typography Scale

#### Font Hierarchy
```css
/* Typography System */
.text-xs { font-size: 0.75rem; line-height: 1rem; }
.text-sm { font-size: 0.875rem; line-height: 1.25rem; }
.text-base { font-size: 1rem; line-height: 1.5rem; }
.text-lg { font-size: 1.125rem; line-height: 1.75rem; }
.text-xl { font-size: 1.25rem; line-height: 1.75rem; }
.text-2xl { font-size: 1.5rem; line-height: 2rem; }
.text-3xl { font-size: 1.875rem; line-height: 2.25rem; }

/* Font Weights */
.font-light { font-weight: 300; }
.font-normal { font-weight: 400; }
.font-medium { font-weight: 500; }
.font-semibold { font-weight: 600; }
.font-bold { font-weight: 700; }
```

### 3. Spacing System

#### Consistent Spacing
```css
/* Spacing Scale */
.space-1 { margin: 0.25rem; }
.space-2 { margin: 0.5rem; }
.space-3 { margin: 0.75rem; }
.space-4 { margin: 1rem; }
.space-6 { margin: 1.5rem; }
.space-8 { margin: 2rem; }
.space-12 { margin: 3rem; }
.space-16 { margin: 4rem; }
```

## Interactive Components

### 1. Button Components

#### Button Variants
```blade
<!-- Primary Button -->
<button class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors">
    Primary Action
</button>

<!-- Secondary Button -->
<button class="px-4 py-2 bg-gray-200 text-gray-800 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors">
    Secondary Action
</button>

<!-- Icon Button -->
<button class="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
    </svg>
    Add Item
</button>
```

### 2. Form Controls

#### Input Fields
```blade
<!-- Text Input -->
<div class="mb-4">
    <label class="block text-sm font-medium text-gray-700 mb-2">
        Email Address
    </label>
    <input 
        type="email" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        placeholder="Enter your email"
    >
</div>

<!-- Select Dropdown -->
<div class="mb-4">
    <label class="block text-sm font-medium text-gray-700 mb-2">
        Category
    </label>
    <select class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
        <option value="">Select a category</option>
        <option value="1">Category 1</option>
        <option value="2">Category 2</option>
    </select>
</div>

<!-- Checkbox -->
<div class="flex items-center mb-4">
    <input 
        type="checkbox" 
        id="terms"
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
    >
    <label for="terms" class="ml-2 block text-sm text-gray-700">
        I agree to the terms and conditions
    </label>
</div>
```

### 3. Navigation Components

#### Navigation Bar
```blade
<nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
            <!-- Logo -->
            <div class="flex items-center">
                <a href="{{ route('home') }}" class="text-xl font-bold text-gray-900">
                    {{ config('app.name') }}
                </a>
            </div>
            
            <!-- Navigation Links -->
            <div class="hidden md:flex items-center space-x-8">
                <a href="{{ route('dashboard') }}" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
                    Dashboard
                </a>
                <a href="{{ route('users.index') }}" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
                    Users
                </a>
                <a href="{{ route('settings') }}" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
                    Settings
                </a>
            </div>
            
            <!-- User Menu -->
            <div class="flex items-center">
                <div class="relative">
                    <button class="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                        <img class="h-8 w-8 rounded-full" src="{{ Auth::user()->avatar_url }}" alt="User avatar">
                    </button>
                </div>
            </div>
        </div>
    </div>
</nav>
```

## Data Display Patterns

### 1. Tables

#### Data Table Component
```blade
<div class="overflow-x-auto">
    <table class="min-w-full bg-white border border-gray-200">
        <thead class="bg-gray-50">
            <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Name
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Email
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                </th>
            </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
            @foreach($users as $user)
                <tr class="hover:bg-gray-50">
                    <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex items-center">
                            <img class="h-10 w-10 rounded-full" src="{{ $user->avatar_url }}" alt="">
                            <div class="ml-4">
                                <div class="text-sm font-medium text-gray-900">{{ $user->name }}</div>
                            </div>
                        </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {{ $user->email }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {{ $user->active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800' }}">
                            {{ $user->active ? 'Active' : 'Inactive' }}
                        </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <a href="{{ route('users.edit', $user) }}" class="text-blue-600 hover:text-blue-900 mr-3">
                            Edit
                        </a>
                        <button onclick="deleteUser({{ $user->id }})" class="text-red-600 hover:text-red-900">
                            Delete
                        </button>
                    </td>
                </tr>
            @endforeach
        </tbody>
    </table>
</div>
```

### 2. Cards

#### Card Component
```blade
<div class="bg-white overflow-hidden shadow-sm rounded-lg">
    <div class="px-6 py-4">
        <h3 class="text-lg font-medium text-gray-900 mb-2">
            {{ $title }}
        </h3>
        <p class="text-sm text-gray-600">
            {{ $description }}
        </p>
    </div>
    @if($actions)
        <div class="px-6 py-3 bg-gray-50 border-t border-gray-200">
            <div class="flex justify-end space-x-3">
                {{ $actions }}
            </div>
        </div>
    @endif
</div>
```

### 3. Status Indicators

#### Status Components
```blade
<!-- Success Status -->
<div class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
    <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
    </svg>
    Active
</div>

<!-- Warning Status -->
<div class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800">
    <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
    </svg>
    Pending
</div>

<!-- Error Status -->
<div class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-800">
    <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
    </svg>
    Failed
</div>
```

## Error Handling and Feedback

### 1. Error Messages

#### Validation Errors
```blade
@if ($errors->any())
    <div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
        <div class="flex">
            <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
                </svg>
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">
                    Please correct the following errors:
                </h3>
                <ul class="mt-2 text-sm text-red-700 list-disc list-inside">
                    @foreach ($errors->all() as $error)
                        <li>{{ $error }}</li>
                    @endforeach
                </ul>
            </div>
        </div>
    </div>
@endif
```

### 2. Success Messages

#### Flash Messages
```blade
@if(session('success'))
    <div class="mb-4 p-4 bg-green-50 border border-green-200 rounded-md">
        <div class="flex">
            <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                </svg>
            </div>
            <div class="ml-3">
                <p class="text-sm font-medium text-green-800">
                    {{ session('success') }}
                </p>
            </div>
        </div>
    </div>
@endif
```

### 3. Loading States

#### Loading Indicators
```blade
<!-- Button Loading State -->
<button 
    type="submit" 
    class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
    :disabled="processing"
>
    <span x-show="!processing">Save Changes</span>
    <span x-show="processing" class="flex items-center">
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Processing...
    </span>
</button>
```

## Responsive Design

### 1. Mobile-First Approach

#### Responsive Classes
```css
/* Mobile First - Base styles */
.container {
    padding: 1rem;
}

/* Tablet - 768px and up */
@media (min-width: 768px) {
    .container {
        padding: 2rem;
    }
}

/* Desktop - 1024px and up */
@media (min-width: 1024px) {
    .container {
        padding: 3rem;
    }
}
```

### 2. Flexible Grid System

#### Grid Layout
```blade
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold mb-2">Card 1</h3>
        <p class="text-gray-600">Content goes here...</p>
    </div>
    <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold mb-2">Card 2</h3>
        <p class="text-gray-600">Content goes here...</p>
    </div>
    <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold mb-2">Card 3</h3>
        <p class="text-gray-600">Content goes here...</p>
    </div>
</div>
```

## Accessibility Guidelines

### 1. ARIA Labels and Roles

#### Accessible Forms
```blade
<form role="form" aria-labelledby="form-title">
    <h2 id="form-title">User Registration</h2>
    
    <div class="mb-4">
        <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email Address
        </label>
        <input 
            type="email" 
            id="email"
            name="email"
            required
            aria-describedby="email-error"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
        <div id="email-error" class="text-red-600 text-sm mt-1" role="alert">
            <!-- Error message will appear here -->
        </div>
    </div>
</form>
```

### 2. Keyboard Navigation

#### Focus Management
```css
/* Focus States */
.btn:focus,
.form-input:focus,
.form-select:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
}

/* Skip Link */
.skip-link {
    position: absolute;
    top: -40px;
    left: 6px;
    background: #000;
    color: #fff;
    padding: 8px;
    text-decoration: none;
    z-index: 1000;
}

.skip-link:focus {
    top: 6px;
}
```

These UI/UX guidelines ensure that {{project}} maintains a consistent, accessible, and professional interface that provides an excellent user experience across all devices and interaction methods.