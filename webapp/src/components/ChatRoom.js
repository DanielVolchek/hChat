import React from "react";
import classes from "../css/chat.module.css";
class ChatRoom extends React.Component {
	constructor(props) {
		super(props);
	}
	render() {
		return (
			<div className={classes.chatModule}>
				<span className={classes.roomTitle}>{this.props.room_desc}</span>
				<div className={classes.groupOfMessages}>{this.props.messages}</div>
				{this.props.scrollDiv}
			</div>
		);
	}
}

export default ChatRoom;
