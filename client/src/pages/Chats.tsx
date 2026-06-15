import { useEffect, useState } from "react";
import { Link } from "react-router";
import { Request, type Chat } from "../conn/network";
import Nav from "../components/Nav";

function ChatDisplay({ chat }: { chat: Chat }) {
	return (
		<Link className="w-1/2 bg-blue-950 text-white rounded-md p-4 m-2" to={`/chat?id=${chat.id}`}>
			<h1 className="text-2xl">{chat.name}</h1>
			<p>{chat.messages.length !== 0 ? chat.messages[chat.messages.length - 1].message : "No messages yet"}</p>
		</Link>
	)
}

function Chats() {
	const [chats, setChats] = useState<Chat[]>([])
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		Request.getChats().then(res => {
			setChats(res.chats);
			setLoading(false);
		}
		).catch(err => console.error(err))
	}, [])

	return (loading) ? (<p>Loading...</p>) : (
		<>
			<Nav />
			<div className="min-h-screen flex justify-center items-center flex-col">
				{(chats.length == 0) ? <h1 className="text-2xl">No chats yet. Add them by clicking + in the bottom right corner.</h1> : chats.map(chat => <ChatDisplay chat={chat} key={chat.id} />)}
			</div>
			<Link to="/createChat" className="absolute w-20 h-20 p-0 m-0 flex items-center justify-center text-4xl bottom-1 right-1 bg-red-950 text-white rounded-full">
				+
			</Link>
		</>
	)
}

export default Chats