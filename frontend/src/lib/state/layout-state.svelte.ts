import { EntityKind } from '$lib/api/canturin';

export type PanelType =
	| 'none'
	| 'network'
	| 'bus'
	| 'node'
	| 'message'
	| 'signal'
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

	openPanel0(entKind: EntityKind, entId: string) {
		let panelType: PanelType;
		switch (entKind) {
			case EntityKind.EntityKindNetwork:
				panelType = 'bus';
				break;
			case EntityKind.EntityKindBus:
				panelType = 'bus';
				break;
			case EntityKind.EntityKindNode:
				panelType = 'node';
				break;
			case EntityKind.EntityKindMessage:
				panelType = 'message';
				break;
			case EntityKind.EntityKindSignal:
				panelType = 'signal';
				break;
			case EntityKind.EntityKindSignalType:
				panelType = 'signal_type';
				break;
			case EntityKind.EntityKindSignalUnit:
				panelType = 'signal_unit';
				break;
			case EntityKind.EntityKindSignalEnum:
				panelType = 'signal_enum';
				break;
			default:
				panelType = 'bus';
		}

		this.openPanel(panelType, entId);
	}
}

export default new LayoutState();
