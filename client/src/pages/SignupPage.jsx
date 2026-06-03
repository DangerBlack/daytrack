import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { generateChallenge as apiChallenge, signup as apiSignup, signin as apiSignin } from '../api/client';
import { generateKeypair, signChallenge, setToken, setUser } from '../api/auth';
import './Auth.css';

export default function SignupPage() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const challenge = await apiChallenge(email);
      const { publicKey, privateKey } = await generateKeypair(password, challenge.salt);
      await apiSignup(username, email, publicKey, challenge.salt);

      const signedChallenge = await signChallenge(privateKey, challenge.challenge);
      const result = await apiSignin(email, challenge.challenge, signedChallenge);
      setToken(result.token);
      setUser({ username, email });
      navigate('/');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="auth-logo">
          <img src="/logo.svg" alt="Daytrack" />
          <h1>Daytrack</h1>
        </div>
        <p className="subtitle">Create your account</p>

        {error && <div className="error">{error}</div>}

        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
          minLength={3}
          pattern="[a-zA-Z0-9_-]+"
          title="Letters, numbers, underscores and hyphens only"
        />
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Creating account...' : 'Sign up'}
        </button>

        <p className="switch">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </form>
    </div>
  );
}
