/// <reference path="../References.d.ts"/>
import * as React from 'react';
import ProfilesStore from '../stores/ProfilesStore';
import * as ProfileTypes from '../types/ProfileTypes';
import * as ProfileActions from '../actions/ProfileActions';
import Profile from "./Profile";

interface State {
	profiles: ProfileTypes.ProfilesRo;
}

const css = {
};

export default class Profiles extends React.Component<{}, State> {
	interval: NodeJS.Timer

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			profiles: ProfilesStore.profiles,
		};
	}

	componentDidMount(): void {
		ProfilesStore.addChangeListener(this.onChange);
		ProfileActions.sync();

		this.interval = setInterval(() => {
			ProfileActions.sync(true)
		}, 1000);
	}

	componentWillUnmount(): void {
		ProfilesStore.removeChangeListener(this.onChange);

		clearInterval(this.interval)
	}

	onChange = (): void => {
		this.setState({
			profiles: ProfilesStore.profiles,
		});
	}

	render(): JSX.Element {
		let profilesDom: JSX.Element[] = [];

		this.state.profiles.forEach((prfl: ProfileTypes.ProfileRo): void => {
			profilesDom.push(<Profile
				key={prfl.id}
				profile={prfl}
			/>);
		});

		return <div>
			{profilesDom}
		</div>;
	}
}
