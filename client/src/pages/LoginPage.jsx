import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { generateChallenge as apiChallenge, signin as apiSignin } from '../api/client';
import { generateKeypair, signChallenge, setToken, setUser } from '../api/auth';
import './Auth.css';

export default function LoginPage() {
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
      const { privateKey } = await generateKeypair(password, challenge.salt);
      const signedChallenge = await signChallenge(privateKey, challenge.challenge);
      const result = await apiSignin(email, challenge.challenge, signedChallenge);
      setToken(result.token);
      setUser({ username: result.username, email });
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
        <p className="subtitle">Sign in to your account</p>

        {error && <div className="error">{error}</div>}

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
          {loading ? 'Signing in...' : 'Sign in'}
        </button>

        <p className="switch">
          Don't have an account? <Link to="/signup">Sign up</Link>
        </p>
      </form>
    </div>
  );
}
