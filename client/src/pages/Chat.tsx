import { useEffect, useState } from "react";
import { Request, type ChatMessage } from "../conn/network";

function ChatBubble({ text, fromUser }: { text: string, fromUser: string }) {
	const userID = localStorage.getItem("userID")!;

	return (fromUser == userID) ? (
		<div className="flex justify-between">
			<div></div>
			<div className="bg-blue-950 w-50 rounded-2xl text-white p-3">
				<p>{text}</p>
			</div>
		</div>
	) : (
		<div className="flex justify-between">
			<div className="bg-red-950 w-50 rounded-2xl text-white p-3">
				<h2>{text}</h2>
				<p>{fromUser}</p>
			</div>
			<div></div>
		</div>
	);
}

function Chat() {
	// const _id = new URLSearchParams(document.location.search);
	const [currMessage, setCurrMessage] = useState("");
	const [messages, setMessages] = useState<ChatMessage[]>([{ author: "afsd", message: "asdf", id: "asdf" }]);
	const [loading, setLoading] = useState(true);
	const userID = localStorage.getItem("userID")!;

	useEffect(() => {
		const params = new URLSearchParams(window.location.search);
		const id = params.get("id");
		if (!id) {
			window.location.href = "/";
			return;
		}
		Request.getMessages({ chat_id: id }).then(res => {
			const mes = res.messages ?? [];
			setMessages(mes);
			setLoading(false);
		})
			.catch(err => console.error(err));
	}, []);

	return (loading) ?
		(<p> Loading...</p >) :
		(<div className="min-h-screen">
			{messages.map((msg, index) => <ChatBubble key={index} text={msg.message} fromUser={msg.author} />)}
			<div className="h-10"></div>
			<form className="absolute bottom-0 w-full flex" onSubmit={(e) => {
				e.preventDefault();
				setMessages([...messages, { id: "nothing", message: currMessage, author: userID }]);
				setCurrMessage("");
			}}>
				<input className="h-10 form-input w-5/6" value={currMessage} onChange={(e) => setCurrMessage(e.target.value)} />
				<input className="h-10 form-input w-1/6" type="submit" value="Send" />
			</form>
		</div >);
}

export default Chat;