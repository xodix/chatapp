import { useEffect, useRef, useState } from "react";
import { Request, type ChatMessage } from "../conn/network";
import Nav from "../components/Nav";

function ChatBubble({ text, fromUser, userName, userSurname }: { text: string; fromUser: string, userName: string, userSurname: string }) {
	const userID = localStorage.getItem("userID")!;
	return fromUser === userID ? (
		<div className="flex justify-between m-1">
			<div></div>
			<div className="bg-blue-950 w-1/3 text-right rounded-2xl text-white p-3">
				<p className="text-xl">{text}</p>
			</div>
		</div>
	) : (
		<div className="flex justify-between m-1">
			<div className="bg-red-950 w-1/3 text-left rounded-2xl text-white p-3">
				<p className="text-xl">{text}</p>
				<p>{userName} {userSurname}</p>
			</div>
			<div></div>
		</div>
	);
}

function Chat() {
	const [currMessage, setCurrMessage] = useState("");
	const [messages, setMessages] = useState<ChatMessage[]>([]);
	const [loading, setLoading] = useState(true);
	const wsRef = useRef<WebSocket | null>(null);

	useEffect(() => {
		const token = localStorage.getItem("token")!;
		const params = new URLSearchParams(window.location.search);
		const id = params.get("id");
		if (!id) {
			window.location.href = "/";
			return;
		}

		Request.getMessages({ chat_id: id })
			.then(res => {
				setMessages(res ?? []);
				setLoading(false);
				setTimeout(() => {
					window.scrollTo(0, document.body.scrollHeight);
				}, 100);
			})
			.catch(err => console.error(err));

		const ws = new WebSocket(`ws://localhost:8080/ws?chat_id=${id}&token=${token}`);
		wsRef.current = ws;

		ws.onopen = () => console.log("connected");
		ws.onmessage = (event) => {
			console.log(event)
			setMessages(prev => [...prev, JSON.parse(event.data)]);
			setTimeout(() => {
				window.scrollTo(0, document.body.scrollHeight);
			}, 100);
		};
		ws.onclose = (event) => console.log("disconnected", event.code);
		ws.onerror = (error) => console.error("error:", error);

		return () => {
			wsRef.current?.close();
			console.log("disconnected");
		}
	}, []);

	const handleSend = (e: React.FormEvent) => {
		e.preventDefault();
		if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN
			|| !currMessage.trim()) return;
		console.log(currMessage, "SENT");
		wsRef.current.send(currMessage);
		setCurrMessage("");
	};

	if (loading) return <p>Loading...</p>;

	return (
		<>
			<Nav />
			<div className="min-h-screen">
				{messages.map((msg, index) => (
					<ChatBubble key={index} text={msg.message} fromUser={msg.author} userName={msg.name} userSurname={msg.surname} />
				))}
				<div className="h-10"></div>
				<form className="fixed bottom-0 w-full flex" onSubmit={handleSend}>
					<input
						className="h-10 form-input w-5/6"
						value={currMessage}
						onChange={(e) => setCurrMessage(e.target.value)}
					/>
					<button className="h-10 w-1/6 form-input" type="submit">Send</button>
				</form>
			</div>
		</>
	);
}

export default Chat;