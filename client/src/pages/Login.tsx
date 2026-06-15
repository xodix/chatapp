import { useState } from 'react';
import { NavLink } from 'react-router';
import { Request } from '../conn/network';

function Login() {
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')

	function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault();

		Request.login({ email, password })
			.then(res => {
				localStorage.setItem('token', res.token);
				localStorage.setItem('userID', res.user_id.toString());
				window.location.href = '/';
			}).catch(err => {
				setError(err.message);
			});
	}

	return (
		<>
			<div className="min-h-screen w-auto flex justify-center items-center">
				<div className="block w-fit">
					<h1 className="text-2xl">Login:</h1>
					<p>{error}</p>
					<form onSubmit={handleSubmit}>
						<label>
							Email:
							<input type="email" className="form-input" name="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
						</label>
						<label>
							Password:
							<input type="password" className="form-input" name="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
						</label>
						<input type="submit" className="form-input w-full" value="Login" />
					</form>
					<p>Or <NavLink to="/register" className="text-blue-500 hover:underline">register</NavLink></p>
				</div>
			</div>
		</>
	)
}

export default Login
