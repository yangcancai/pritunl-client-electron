/// <reference path="../References.d.ts"/>
import * as React from "react"
import * as Theme from "../Theme"
import ProfilesStore from "../stores/ProfilesStore"
import * as ProfileTypes from "../types/ProfileTypes"
import * as ProfileActions from "../actions/ProfileActions"
import * as ServiceActions from "../actions/ServiceActions"
import * as Blueprint from "@blueprintjs/core"
import PageInfo from "./PageInfo"
import PageSwitch from "./PageSwitch"
import AceEditor from "react-ace"

import "ace-builds/src-noconflict/mode-text"
import "ace-builds/src-noconflict/theme-dracula"
import "ace-builds/src-noconflict/theme-eclipse"
import * as MiscUtils from "../utils/MiscUtils"
import * as Constants from "../Constants"

interface Props {
	profile: ProfileTypes.ProfileRo
	onConfirm?: () => void
}

interface State {
	disabled: boolean
	password: string
	dialog: boolean
	confirm: number
	confirming: string
}

const css = {
	box: {
		display: "inline-block"
	} as React.CSSProperties,
	button: {
		marginRight: "10px",
	} as React.CSSProperties,
	dialog: {
		width: "340px",
		position: "absolute",
	} as React.CSSProperties,
	label: {
		width: "100%",
		maxWidth: "220px",
		margin: "18px 0 0 0",
	} as React.CSSProperties,
	input: {
		width: "100%",
	} as React.CSSProperties,
}

export default class ProfileConnect extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context)
		this.state = {
			disabled: false,
			password: "",
			dialog: false,
			confirm: 0,
			confirming: null,
		}
	}

	openDialog = (): void => {
		this.setState({
			...this.state,
			dialog: true,
		})
	}

	closeDialog = (): void => {
		this.setState({
			...this.state,
			dialog: false,
		})
	}

	closeDialogConfirm = (): void => {
		this.setState({
			...this.state,
			dialog: false,
		})
		if (this.props.onConfirm) {
			this.props.onConfirm()
		}
	}

	clearConfirm = (): void => {
		this.setState({
			...this.state,
			confirm: 0,
			confirming: null,
		})
	}

	render(): JSX.Element {
		let buttonClass = ""
		let buttonLabel = ""
		if (this.props.profile.status === "connected") {
			buttonClass = "bp3-intent-danger bp3-icon-link"
			buttonLabel = "Disconnect"
		} else {
			buttonClass = "bp3-intent-success bp3-icon-link"
			buttonLabel = "Connect"
		}

		return <div style={css.box}>
			<button
				className={"bp3-button " + buttonClass}
				style={css.button}
				type="button"
				disabled={this.state.disabled}
				onClick={this.openDialog}
			>
				{buttonLabel}
			</button>
			<Blueprint.Dialog
				title="Profile Connect"
				style={css.dialog}
				isOpen={this.state.dialog}
				usePortal={true}
				portalContainer={document.body}
				onClose={this.closeDialog}
			>
				<div className="bp3-dialog-body">
					Connecting to {this.props.profile.formattedName()}
					<label
						className="bp3-label"
						style={css.label}
					>
						Enter password:
						<input
							className="bp3-input"
							style={css.input}
							disabled={this.state.disabled}
							autoCapitalize="off"
							spellCheck={false}
							placeholder="Enter password"
							value={this.state.password}
							onChange={(evt): void => {
								this.setState({
									...this.state,
									password: evt.target.value,
								})
							}}
						/>
					</label>
				</div>
				<div className="bp3-dialog-footer">
					<div className="bp3-dialog-footer-actions">
						<button
							className="bp3-button bp3-intent-danger"
							type="button"
							onClick={this.closeDialog}
						>Cancel</button>
						<button
							className="bp3-button bp3-intent-success bp3-icon-link"
							type="button"
							disabled={this.state.disabled || this.state.password === ""}
							onClick={this.closeDialogConfirm}
						>Connect</button>
					</div>
				</div>
			</Blueprint.Dialog>
		</div>
	}
}