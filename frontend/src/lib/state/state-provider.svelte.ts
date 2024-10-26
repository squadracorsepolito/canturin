import type { WailsEvent } from '@wailsio/runtime/types/events';
import type { Entity, EntityState } from './entity-state.svelte';
import * as wails from '@wailsio/runtime';

export class StateProvider<E extends Entity, S extends EntityState<E>> {
	states = $state(new Map<string, S>());
	genFn: (entity: E) => S;

	constructor(stateFactory: (entity: E) => S, historyEventName: string) {
		this.genFn = stateFactory;

		wails.Events.On(historyEventName, (e: WailsEvent) => {
			this.modify(e.data[0] as E);
		});
	}

	add(entity: E) {
		this.states.set(entity.entityId, this.genFn(entity));
	}

	remove(entityId: string) {
		this.states.delete(entityId);
	}

	get(entityId: string) {
		const data = this.states.get(entityId);
		if (!data) {
			throw new Error(`No entity with id ${entityId} found`);
		}
		return data as S;
	}

	modify(entity: E) {
		const s = this.states.get(entity.entityId);
		if (s) {
			s.entity = entity;
		}
	}
}
