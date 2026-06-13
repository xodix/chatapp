import { useEffect, useState } from "react";
import { Request, type User } from "../conn/network";


function CreateChat() {
	const [name, setName] = useState('')
	const [loadingUsers, setLoadingUsers] = useState(true)
	const [allMembers, setAllMembers] = useState<User[]>([]);
	const [members, setMembers] = useState<string[]>([]);
	const [error, setError] = useState('')
	const userID = localStorage.getItem('userID')!;

	useEffect(() => {
		const userID = localStorage.getItem('userID')!;
		Request.getAllUsers().then(res => {
			setAllMembers(res.users.filter(user => user.id != userID));
			setLoadingUsers(false);
		}).catch(err => {
			setError(err.message);
		});
	}, []);

	function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault();

		Request.createChat({ name, members })
			.then(res => {
				console.log(res);
				location.href = '/';
			}).catch(err => {
				setError(err.message);
			});
	}

	return (
		<>
			<div className="min-h-screen w-auto flex justify-center items-center">
				<div className="block w-fit">
					<h1 className="text-2xl">Create a chat:</h1>
					<p>{error}</p>
					<form onSubmit={handleSubmit}>
						<label>
							Chat name:
							<input type="text" className="form-input" name="name" placeholder="name" value={name} onChange={(e) => setName(e.target.value)} />
						</label>
						<label>
							Members:
							<select name="members" className="form-input" multiple onChange={(e) => setMembers([...[...e.target.options].filter(option => option.selected).map(option => option.value), userID])}>
								{loadingUsers ? <option disabled>Loading users...</option> : allMembers.map(user => {
									return <option key={user.id} value={user.id}>{user.name} {user.surname}</option>
								})}
							</select>
						</label>
						<input type="submit" className="form-input w-full" value="Create" />
					</form>
				</div>
			</div>
		</>
	);
}

export default CreateChat;