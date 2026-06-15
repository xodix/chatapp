import Nav from "../components/Nav"

export default function NotFound() {
	return (
		<>
			<Nav />
			<div className="min-h-screen flex justify-center items-center flex-col">
				<h1 className="text-2xl">Could not find the page you were looking for.</h1>
			</div>
		</>
	)
}