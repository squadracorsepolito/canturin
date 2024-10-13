type PanelType = 'none' | 'message' | 'signal_type' | 'signal_unit';

class LayoutState {
	openPanelType: PanelType = $state('none');
	openPanelId = $state('');

	async openPanel(typ: PanelType, panelId: string) {
		this.openPanelType = typ;
		this.openPanelId = panelId;
	}
}

export default new LayoutState();
