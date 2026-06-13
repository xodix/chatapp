import { useState } from 'react'
import { NavLink } from 'react-router'
import { Request } from '../conn/network'

function leftPad(text: string, totalLength: number, fillChar: string): string {
	let buff = "";
	for (let i = 0; i < totalLength - text.length; i++) {
		buff += fillChar;
	}
	buff += text;

	return buff;
}

function formatDate(date: Date): string {
	const yyyy = leftPad(date.getFullYear().toString(), 4, "0");
	const mm = leftPad((date.getMonth() + 1).toString(), 2, "0");
	const dd = leftPad(date.getDate().toString(), 2, "0");

	return `${yyyy}-${mm}-${dd}`;
}

function Register() {
	const now = new Date();
	const minimumDate = new Date(
		now.getFullYear() - 18,
		now.getMonth(),
		now.getDate()
	);

	const [name, setName] = useState('')
	const [surname, setSurname] = useState('')
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [birthdate, setBirthdate] = useState(new Date(
		minimumDate.getFullYear(),
		0,
		1
	));
	const [error, setError] = useState('')

	function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault();

		Request.register({ name, surname, email, password, birthdate }).then(res => {
			setError("");
			localStorage.setItem('token', res.token);
			localStorage.setItem('userID', res.user_id.toString());
			window.location.href = '/';
		}).catch(err => {
			setError(err.message);
		})
	}

	return (
		<div className="min-h-screen w-auto flex justify-center items-center">
			<div className="block w-fit">
				<h1 className="text-2xl">Register</h1>
				{error ? <p className="text-red-500 font-bold">{error}</p> : <p></p>}
				<form onSubmit={handleSubmit}>
					<label>
						name:
						<input type="text" className="form-input" name='name' placeholder='Name' value={name} onChange={(e) => setName(e.target.value)} />
					</label>
					<label>
						surname:
						<input type="text" className="form-input" name='surname' placeholder='Surname' value={surname} onChange={(e) => setSurname(e.target.value)} />
					</label>
					<label>
						email:
						<input type="email" className="form-input" name="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
					</label>
					<label>
						password:
						<input type="password" className="form-input" name="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
					</label>
					<label>
						Birthdate:
						<input
							type="date"
							max={formatDate(minimumDate)}
							className="form-input w-full"
							name="birthdate"
							placeholder="Birthdate"
							value={formatDate(birthdate)}
							onChange={(e) => setBirthdate(new Date(e.target.value))} />
					</label>
					<input type="submit" className="form-input w-full" value="Register" />
				</form>
				<p>Or <NavLink to="/login" className="text-blue-500 hover:underline">log in</NavLink></p>
			</div>
		</div>
	)
}

export default Register
