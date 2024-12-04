export type PanelType =
	| 'none'
	| 'bus'
	| 'node'
	| 'message'
	| 'signal_type'
	| 'signal_unit'
	| 'signal_enum';

class LayoutState {
	openPanelType: PanelType = $state('none');
	openPanelId = $state('');

	openPanel(typ: PanelType, panelId: string) {
		this.openPanelType = typ;
		this.openPanelId = panelId;
	}

	openMessagePanel(msgEntId: string) {
		this.openPanel('message', msgEntId);
	}

	openMessageDraftPanel() {
		this.openPanel('message', 'draft');
	}

	closeIfOpen(panelId: string) {
		if (this.openPanelId === panelId) {
			this.openPanel('none', '');
		}
	}
}

export default new LayoutState();
