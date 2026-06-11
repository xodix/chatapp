import { useState } from "react";

function ChatBubble({ text, fromUser }: { text: string, fromUser: boolean }) {
	return (fromUser) ? (
		<div className="flex justify-between">
			<div></div>
			<div className="bg-blue-950 w-50 rounded-2xl text-white p-3">
				<p>{text}</p>
			</div>
		</div>
	) : (
		<div className="flex justify-between">
			<div className="bg-red-950 w-50 rounded-2xl text-white p-3">
				<p>{text}</p>
			</div>
			<div></div>
		</div>
	);
}

function Chat() {
	// const _id = new URLSearchParams(document.location.search);
	const [currMessage, setCurrMessage] = useState("");
	const messages = [
		{ text: "Hello", fromUser: true },
		{ text: "Mellow my friend", fromUser: false },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
		{ text: "Fuck you", fromUser: true },
	];

	return (
		<div className="min-h-screen">
			{messages.map((msg, index) => (
				<ChatBubble key={index} text={msg.text} fromUser={msg.fromUser} />
			))}
			<div className="h-10"></div>
			<form className="absolute bottom-0 w-full flex" onSubmit={(e) => {
				e.preventDefault();
				setCurrMessage("");
			}}>
				<input className="h-10 form-input w-5/6" value={currMessage} onChange={(e) => setCurrMessage(e.target.value)} />
				<input className="h-10 form-input w-1/6" type="submit" value="Send" />
			</form>
		</div>
	);
}

export default Chat;