<script lang="ts">
    import { goto } from '$app/navigation';
    
    // Form State (using Svelte's reactive declarations)
    let formData = {
        email: "",
        password: "",
    };
    
    // UI State
    let isLoading = false;
    let error: string | null = null;

    // TypeScript Interface for API Response
    interface LoginResponse {
        message: string;
        token?: string;
    }

    // Function to handle form submission and API call
    const handleSubmit = async () => {
        isLoading = true;
        error = null;

        try {
            const response = await fetch("http://localhost:8088/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(formData),
            });

            const data: LoginResponse = await response.json();

            if (response.ok) {
                // SUCCESS: Save token (e.g., to a cookie) and redirect
                console.log("Login successful:", data.token);
                goto('/dashboard'); 
            } else {
                // FAILURE
                error = data.message || "Login failed. Please check your credentials.";
            }
        } catch (err) {
            console.error("Network or server error:", err);
            error = "Could not connect to the server. Please try again later.";
        } finally {
            isLoading = false;
        }
    };
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
    <div class="w-full max-w-md">
        <div class="bg-white rounded-lg shadow-lg p-8">
            <div class="flex items-center justify-center mb-8">
                <div class="w-12 h-12 bg-gradient-to-br from-blue-600 to-blue-400 rounded-lg flex items-center justify-center">
                    <span class="text-white font-bold text-xl">U</span>
                </div>
                <h1 class="text-2xl font-bold text-gray-900 ml-3">Unimart</h1>
            </div>

            <h2 class="text-2xl font-bold text-gray-900 mb-2">Welcome Back</h2>
            <p class="text-gray-600 mb-6">Sign in to your account to continue.</p>

            <form on:submit|preventDefault={handleSubmit} class="space-y-4">
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700 mb-2">University Email</label>
                    <input
                        id="email"
                        type="email"
                        name="email"
                        bind:value={formData.email}
                        placeholder="your.name@example.com"
                        required
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    />
                </div>

                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700 mb-2">Password</label>
                    <input
                        id="password"
                        type="password"
                        name="password"
                        bind:value={formData.password}
                        placeholder="Enter your password"
                        required
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    />
                </div>

                {#if error}
                    <p class="text-sm text-red-600 font-medium text-center">{error}</p>
                {/if}

                <button
                    type="submit"
                    disabled={isLoading}
                    class="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-2 rounded-lg hover:shadow-lg transition-shadow disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {#if isLoading}
                        Signing In...
                    {:else}
                        Sign In
                    {/if}
                </button>
            </form>

            <p class="text-center text-gray-600 mt-4">
                Don't have an account? 
                <a href="/signup" class="text-blue-600 font-semibold cursor-pointer hover:underline">Sign Up</a>
            </p>
        </div>
    </div>
</div>
