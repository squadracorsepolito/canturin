import { getContext, setContext } from 'svelte';

type PanelType = 'none' | 'message' | 'signal_type' | 'signal_unit';

class LayoutState {
	openPanelType: PanelType = $state('none');
	openPanelId = $state('');

	async openPanel(typ: PanelType, panelId: string) {
		this.openPanelType = typ;
		this.openPanelId = panelId;
	}
}

const LAYOUT_KEY = Symbol('LAYOUT');

export function setLayoutState() {
	return setContext(LAYOUT_KEY, new LayoutState());
}

export function getLayoutState() {
	return getContext<ReturnType<typeof setLayoutState>>(LAYOUT_KEY);
}
