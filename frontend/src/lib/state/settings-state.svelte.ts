import { SettingsService, type Settings } from '$lib/api/canturin';

export let state: SettingsState;

export async function loadSettings() {
	const settings = await SettingsService.Get();
	state = new SettingsState(settings);
}

class SettingsState {
	settings = $state() as Settings;

	constructor(settings: Settings) {
		this.settings = settings;
	}
}
