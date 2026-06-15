import { Link } from "react-router-dom";

export default function Nav() {
	const authorized = !!localStorage.getItem('token');

	if (authorized) return (
		<>
			<nav className="bg-white border-b border-gray-100 dark:bg-gray-900 dark:border-gray-800 flex justify-between items-center fixed w-full top-0 left-0 z-50 text-white p-2">
				<Link to="/"><h1 className="text-2xl">Chats</h1></Link>
				<button className="cursor-pointer" onClick={() => { localStorage.clear(); window.location.href = '/login'; }}><h1 className="text-2xl">Log out</h1></button>
			</nav>

			<div className="h-16"></div>
		</>
	);

	return (
		<>
			<nav className="bg-white border-b border-gray-100 dark:bg-gray-900 dark:border-gray-800 flex justify-between items-center fixed w-full top-0 left-0 z-50 text-white p-2">
				<Link to="/login"><h1 className="text-2xl">Login</h1></Link>
				<Link to="/register"><h1 className="text-2xl">Register</h1></Link>
			</nav>

			<div className="h-16"></div>
		</>
	);

}