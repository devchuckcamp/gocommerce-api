# Google OAuth Setup and Usage

## Configuration

The Google OAuth is configured via environment variables in `.env`:

```env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

## Google OAuth Flow

### 1. Get Authorization URL

**Endpoint:** `GET /api/v1/auth/google`

**Optional Query Parameters:**
- `state` - CSRF protection token (auto-generated if not provided)

**Response:**
```json
{
  "data": {
    "url": "https://accounts.google.com/o/oauth2/auth?...",
    "state": "random-state-..."
  }
}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/auth/google
```

### 2. Redirect User to Google

Direct the user to the `url` returned from step 1. They will:
1. Log in to their Google account (if not already logged in)
2. Grant permissions to your application
3. Be redirected back to your callback URL

### 3. Handle Callback

**Endpoint:** `GET /api/v1/auth/google/callback`

**Query Parameters:**
- `code` - Authorization code from Google
- `state` - CSRF token (should match the one from step 1)

**Response:**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@gmail.com",
      "first_name": "John",
      "last_name": "Doe",
      "email_verified": true,
      "active": true
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "random-token-string",
    "expires_at": "2025-11-29T12:00:00Z"
  }
}
```

## OAuth Behavior

### New User
When a user logs in with Google for the first time:
1. A new user account is created with their Google profile information
2. An OAuth account record is linked to the user
3. Email is automatically verified (if verified by Google)
4. No password is set (OAuth-only authentication)

### Existing User (Same Email)
If a user already exists with the same email:
1. The Google OAuth account is linked to the existing user
2. User can now log in with either password or Google OAuth
3. Email verification status is updated if verified by Google

### Returning OAuth User
When a user logs in with Google again:
1. OAuth tokens are updated
2. User profile information is refreshed from Google
3. New JWT tokens are issued

## Testing with Browser

1. **Get the authorization URL:**
   ```bash
   curl http://localhost:8080/api/v1/auth/google
   ```

2. **Copy the `url` from the response and open it in your browser**

3. **After authorizing, Google will redirect you to:**
   ```
   http://localhost:8080/api/v1/auth/google/callback?code=...&state=...
   ```

4. **The response will contain your user info and JWT tokens**

## Integration in Frontend

### React/Next.js Example

```javascript
// 1. Get OAuth URL
const getGoogleAuthUrl = async () => {
  const response = await fetch('http://localhost:8080/api/v1/auth/google');
  const data = await response.json();
  return data.data;
};

// 2. Redirect to Google
const handleGoogleLogin = async () => {
  const { url, state } = await getGoogleAuthUrl();
  // Store state in session storage for verification
  sessionStorage.setItem('oauth_state', state);
  // Redirect to Google
  window.location.href = url;
};

// 3. Handle callback (on your callback page)
const handleCallback = async () => {
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  const state = params.get('state');
  
  // Verify state matches
  const savedState = sessionStorage.getItem('oauth_state');
  if (state !== savedState) {
    console.error('Invalid state parameter');
    return;
  }
  
  // Exchange code for tokens (backend handles this automatically)
  // The response from the callback endpoint contains the user and tokens
  const response = await fetch(window.location.href);
  const data = await response.json();
  
  // Store tokens
  localStorage.setItem('access_token', data.data.access_token);
  localStorage.setItem('refresh_token', data.data.refresh_token);
  
  // Redirect to dashboard
  window.location.href = '/dashboard';
};
```

## Security Considerations

1. **CSRF Protection:** Always verify the `state` parameter matches what you sent
2. **HTTPS in Production:** Use HTTPS for the redirect URL in production
3. **Store Tokens Securely:** Keep JWT tokens in httpOnly cookies or secure storage
4. **Validate Tokens:** Always validate JWT tokens on protected endpoints

## Google Cloud Console Setup

Make sure your redirect URL is registered in Google Cloud Console:
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to APIs & Services > Credentials
3. Select your OAuth 2.0 Client ID
4. Add `http://localhost:8080/api/v1/auth/google/callback` to Authorized redirect URIs
5. For production, add your production domain

## Troubleshooting

### "redirect_uri_mismatch" Error
- Ensure the redirect URL in `.env` matches exactly what's configured in Google Cloud Console
- Check for trailing slashes and protocol (http vs https)

### "invalid_client" Error
- Verify your Client ID and Client Secret are correct
- Ensure OAuth is enabled in your configuration

### User Not Created
- Check server logs for detailed error messages
- Verify database connection and migrations are complete
