import { useState } from 'react'
import { NavLink } from 'react-router'

function Register() {
	const [username, setUsername] = useState('')
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')

	function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault();

		fetch('http://localhost:5000/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ username, email, password }),
		})
			.then(res => {
				if (res.ok) {
					return res.json();
				}
			}).catch(err => {
				setError(err.message);
			});
	}

	return (
		<div className="min-h-screen w-auto flex justify-center items-center">
			<div className="block w-fit">
				<h1 className="text-2xl">Register</h1>
				<p>{error}</p>
				<form onSubmit={handleSubmit}>
					<label>
						username:
						<input type="text" className="form-input" name='username' placeholder='Username' value={username} onChange={(e) => setUsername(e.target.value)} />
					</label>
					<label>
						email:
						<input type="email" className="form-input" name="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
					</label>
					<label>
						password:
						<input type="password" className="form-input" name="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
					</label>
					<input type="submit" className="form-input w-full" value="Register" />
				</form>
				<p>Or <NavLink to="/login" className="text-blue-500 hover:underline">log in</NavLink></p>
			</div>
		</div>
	)
}

export default Register
