import { ConfigService, type Config } from '$lib/api/canturin';

export let state: ConfigState;

export async function loadConfig() {
	const cfg = await ConfigService.Get();
	state = new ConfigState(cfg);
}

class ConfigState {
	cfg = $state() as Config;

	constructor(cfg: Config) {
		this.cfg = cfg;
	}
}
