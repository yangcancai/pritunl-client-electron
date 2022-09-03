/// <reference path="../References.d.ts"/>
import * as React from 'react'
import ProfilesStore from '../stores/ProfilesStore'
import * as ProfileTypes from '../types/ProfileTypes'
import * as ProfileActions from '../actions/ProfileActions'
import * as Theme from "../Theme";
import * as Constants from "../Constants"
import * as LogUtils from "../utils/LogUtils"
import ConfirmButton from "./ConfirmButton"

import AceEditor from "react-ace"
import "ace-builds/src-noconflict/mode-text"
import "ace-builds/src-noconflict/theme-dracula"
import "ace-builds/src-noconflict/theme-eclipse"
import {Ace} from "ace-builds"

interface State {
	profiles: ProfileTypes.ProfilesRo
	curProfile: ProfileTypes.Profile
	view: string
	log: string
	disabled: boolean
}

const css = {
	message: {
		margin: '0 0 6px 0',
	} as React.CSSProperties,
	header: {
		margin: '0 0 5px 0',
	} as React.CSSProperties,
	label: {
		marginBottom: '0',
	} as React.CSSProperties,
	card: {
		position: "relative",
		margin: '8px',
		paddingRight: 0,
	} as React.CSSProperties,
	layout: {
		height: '100%',
	} as React.CSSProperties,
	progress: {
		width: '100%',
		height: '4px',
		borderRadius: 0,
	} as React.CSSProperties,
	progressBar: {
		width: '50%',
		borderRadius: 0,
	} as React.CSSProperties,
	button: {
		marginRight: '10px',
	} as React.CSSProperties,
	deleteButton: {
	} as React.CSSProperties,
	deleteButtonBox: {
		position: "absolute",
		top: "5px",
		right: "5px",
	} as React.CSSProperties,
	buttons: {
		flexShrink: 0,
	} as React.CSSProperties,
	editor: {
		margin: '10px 20px 0 0',
	} as React.CSSProperties,
};

export default class Logs extends React.Component<{}, State> {
	editor: Ace.Editor

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			profiles: ProfilesStore.profiles,
			curProfile: null,
			view: "service",
			log: "",
			disabled: false,
		};
	}

	componentDidMount(): void {
		Constants.addChangeListener(this.onChange)
		ProfilesStore.addChangeListener(this.onChange);
		ProfileActions.sync();
		this.onChange()
	}

	componentWillUnmount(): void {
		Constants.removeChangeListener(this.onChange)
		ProfilesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		if (this.state.view === "service") {
			LogUtils.readServiceLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					profiles: ProfilesStore.profiles,
				})
			})
		} else if (this.state.view === "client") {
			LogUtils.readClientLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					profiles: ProfilesStore.profiles,
				})
			})
		} else if (this.state.curProfile) {
			this.state.curProfile.readLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					profiles: ProfilesStore.profiles,
				})
			})
		}
	}

	onChangeView = (view: string): void => {
		if (view === "service") {
			LogUtils.readServiceLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					view: view,
					profiles: ProfilesStore.profiles,
				})
			})
		} else if (view === "client") {
			LogUtils.readClientLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					view: view,
					profiles: ProfilesStore.profiles,
				})
			})
		} else {
			let prfl = ProfilesStore.profile(view)

			prfl.readLog().then((data: string) => {
				this.setState({
					...this.state,
					log: data,
					view: view,
					curProfile: prfl,
					profiles: ProfilesStore.profiles,
				})
			})
		}
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		})

		if (this.state.view === "service") {
			LogUtils.clearServiceLog().then((): void => {
				this.setState({
					...this.state,
					disabled: false,
				})
				this.onChange()
			})
		} else if (this.state.view === "client"){
			LogUtils.clearClientLog().then((): void => {
				this.setState({
					...this.state,
					disabled: false,
				})
				this.onChange()
			})
		} else if (this.state.curProfile) {
			this.state.curProfile.clearLog().then((): void => {
				this.setState({
					...this.state,
					disabled: false,
				})
				ProfileActions.sync()
			})
		}
	}

	render(): JSX.Element {
		let label = ""
		if (this.state.view === "service") {
			label = "Service"
		} else if (this.state.view === "client") {
			label = "Client"
		} else if (this.state.curProfile) {
			label = this.state.curProfile.formattedName()
		}

		let viewsDom: JSX.Element[] = [
			<option value="service">Service logs</option>,
			<option value="client">Client logs</option>,
		]

		this.state.profiles.forEach((prfl: ProfileTypes.ProfileRo): void => {
			viewsDom.push(<option
				value={prfl.id}
			>{prfl.formattedName() + " logs"}</option>)
		})

		return <div className="bp3-card layout vertical flex" style={css.card}>
			<div style={css.deleteButtonBox}>
				<ConfirmButton
					className="bp3-minimal bp3-intent-danger bp3-icon-trash"
					style={css.deleteButton}
					safe={true}
					progressClassName="bp3-intent-danger"
					dialogClassName="bp3-intent-danger bp3-icon-delete"
					dialogLabel="Clear Logs"
					confirmMsg={"Confirm clearing " + label + " logs"}
					disabled={this.state.disabled}
					onConfirm={this.onDelete}
				/>
			</div>
			<div className="layout horizontal">
				<h3 style={css.header}>Log Viewer</h3>
			</div>
			<div className="layout horizontal">
				<div className="bp3-select">
					<select
						disabled={this.state.disabled}
						value={this.state.view}
						onChange={(evt): void => {
							this.onChangeView(evt.target.value)
						}}
					>
						{viewsDom}
					</select>
				</div>
			</div>
			<div className="layout horizontal flex">
				<label
					className="bp3-label flex"
					style={css.editor}
				>
					<AceEditor
						name="log-view"
						theme={Theme.editorTheme()}
						height="100%"
						width="100%"
						mode="text"
						fontSize="10px"
						wrapEnabled={true}
						showPrintMargin={false}
						showGutter={true}
						readOnly={true}
						value={this.state.log}
						editorProps={{
							$blockScrolling: true,
						}}
						setOptions={{
							showFoldWidgets: false,
						}}
						onLoad={(editor: Ace.Editor): void => {
							this.editor = editor
							this.editor.scrollToLine(Number.POSITIVE_INFINITY,
								false, false, null)
						}}
					/>
				</label>
			</div>
		</div>;
	}
}
