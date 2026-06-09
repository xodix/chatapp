import { Link } from "react-router";

class Chat {
	uuid: string;
	name: string;
	lastMessage: string;

	constructor(uuid: string, name: string, lastMessage: string) {
		this.uuid = uuid;
		this.name = name;
		this.lastMessage = lastMessage;
	}
}

function ChatDisplay(chat: Chat) {
	return (
		<Link className="w-1/2 bg-blue-950 text-white rounded-md p-4 m-2" to={`/chat?id=${chat.uuid}`}>
			<h1 className="text-2xl">{chat.name}</h1>
			<p>{chat.lastMessage}</p>
		</Link>
	)
}

function Chats() {
	const sampleChats: Chat[] = [
		new Chat("c313a49e-2f48-438a-9199-a8ea6c967c2a", "idiot center", "based"),
		new Chat("c313a49e-2f48-438a-9199-a8ea6c967c2a", "idiot center", "based"),
		new Chat("c313a49e-2f48-438a-9199-a8ea6c967c2a", "idiot center", "based"),
		new Chat("c313a49e-2f48-438a-9199-a8ea6c967c2a", "idiot center", "based"),
	]

	return (
		<div className="min-h-screen flex justify-center items-center flex-col">
			{sampleChats.map(chat => ChatDisplay(chat))}
		</div>
	)
}

export default Chats