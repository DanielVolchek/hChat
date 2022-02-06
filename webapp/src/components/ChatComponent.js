import React from "react";
import Message from "./Message";
import ChatRoom from "./ChatRoom";
import MessageInput from "./MessageInput";
import classes from "../css/chat.module.css";

class ChatComponent extends React.Component {
	constructor(props) {
		// connect to server here
		super(props);
		this.ws = null;
		this.state = {
			messages: [],
			inputvalue: "",
		};
		//Div to scroll to bottom of screen
		this.scrollDiv = React.createRef() 
		// function binds
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
		this.sendToServer = this.sendToServer.bind(this);
		this.receiveFromServer = this.receiveFromServer.bind(this);
	}


	connectToServer() {
		// TODO
		this.ws.onopen = () => {
			this.ws.send(
				`#${this.props.room.room_num}_${this.props.username}:`
			);
		};
		this.ws.onmessage = (msg) => {
			this.receiveFromServer(msg);
		};
	}
	sendToServer(msg) {
		console.log(`Sending ${msg} to server`)
		if (this.validateMessage(msg)) {
			this.ws.send(this.processMessage(msg));
		}
	}
	receiveFromServer(rcv_msg) {
		// TODO
		// Parse received message
		// And make new message component
		// with time and username and message
		const data = rcv_msg.data;
		console.log(this);
		let user = null;
		let msg = null;
		for (let i = 0; i < data.length; i++) {
			const user_end_index = data.indexOf(":");
			user = data.substring(data.indexOf("_") + 1, user_end_index);
			msg = data.substring(user_end_index + 1);
		}
		console.log(`received message ${msg} from server`);
		if (msg === "END_OF_DB") {
			const length = this.state.messages.length
			console.log("END_OF_DB")
			console.log(this.scrollDiv)
			this.scrollDiv.current.scrollIntoView({ behavior: "smooth" })
		}
		const time = new Date();
		if (this.isValidChat(rcv_msg)) {
			const message = <Message time={time} username={user} message={msg} />;
			this.setState({
				...this.state,
				messages: this.state.messages.concat(message),
			})
		};
	}

	isValidChat(message) {
		// TODO
		// placeholder
		return true
	}
	validateMessage(msg) {
		// placeholder
		return msg.length > 0 && msg.length <= 500;
	}
	processMessage(msg) {
		return `#${this.props.room.room_num}_${this.props.username}:${msg}`;
	}
	handleSubmit(event) {
		event.preventDefault();
		this.sendToServer(event.target[0].value);
		this.setState({
			...this.state,
			inputvalue: "",
		});
		// Send message to server
	}
	handleChange(event) {
		this.setState({
			...this.state,
			inputvalue: event.target.value,
		});
	}
	attemptConnection() {
		// If ws connection already exists or username is empty
		if (this.ws !== null || this.props.username === "") {
			return;
		}
		this.ws = new WebSocket("ws://hchat.org:5000/ws");
		this.connectToServer();
	}
	componentDidUpdate(prevProps) {
		if (this.props.username !== prevProps.username) {
			this.attemptConnection();
		}
		if (this.props.room !== prevProps.room) {
			this.setState({
				...this.state,
				messages: [],
			});
			this.ws.send(this.processMessage(`${this.props.username} joined the room`));
		}
	}
	render() {
		return (
			<div className={classes.chatComponent}>
				<ChatRoom
					room_desc={this.props.room.room_desc}
					messages={this.state.messages}
					scrollDiv={
						<div
						ref={ this.scrollDiv }
						></div>
					}
				/>
				<MessageInput
					handleSubmit={this.handleSubmit}
					handleChange={this.handleChange}
					value={this.state.inputvalue}
					room_desc={this.props.room.room_desc}
				/>
			</div>
		);
	}
}

export default ChatComponent;
