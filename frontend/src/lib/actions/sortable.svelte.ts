import { combine } from '@atlaskit/pragmatic-drag-and-drop/combine';
import {
	draggable,
	dropTargetForElements,
	monitorForElements
} from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
import {
	attachClosestEdge,
	extractClosestEdge
} from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
import type { Edge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/dist/types/types';
import { getReorderDestinationIndex } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/get-reorder-destination-index';
import { onMount } from 'svelte';

const privateItemKey = Symbol('Item');

type Item = {
	[privateItemKey]: true;
	instanceId: string;
	id: string;
};

function getItem(data: Omit<Item, typeof privateItemKey>) {
	return {
		[privateItemKey]: true,
		...data
	};
}

function isItem(data: Record<string | symbol, unknown>): data is Item {
	return Boolean(data[privateItemKey]);
}

type Options<T extends { id: string }> = {
	instanceId: string;
	enabled?: boolean;
	itemsGetter: () => T[];
	reorder: (id: string, from: number, to: number) => void;
};

type PartOptions = {
	id: string;
};

export class Sortable<T extends { id: string }> {
	#instanceId: string;
	#itemsGetter: () => T[];
	#reorder: (id: string, from: number, to: number) => void;

	#dropIndicators = new Map<string, { el: HTMLElement; edge: Edge | null }>();

	#dragHandles = new Map<string, HTMLElement>();

	selectedItem = $state({
		id: '',
		index: -1
	});

	mode = $state<'drag' | 'keyboard'>('drag');

	enabled = $state(false);

	constructor({ instanceId, enabled, itemsGetter, reorder }: Options<T>) {
		this.#instanceId = instanceId;
		this.#itemsGetter = itemsGetter;
		this.#reorder = reorder;

		if (enabled) {
			this.enabled = enabled;
		}

		this.mount();
	}

	private mount() {
		const instanceId = this.#instanceId;
		const itemsGetter = this.#itemsGetter;
		const reorder = this.#reorder;

		onMount(() => {
			return monitorForElements({
				canMonitor: ({ source }) => {
					return isItem(source.data) && source.data.instanceId === instanceId && this.enabled;
				},
				onDrop: ({ source, location }) => {
					if (location.current.dropTargets.length === 0) return;

					if (!isItem(source.data)) return;

					const itemId = source.data.id;

					const itemIdx = itemsGetter().findIndex((item) => item.id === itemId);
					if (itemIdx === -1) return;

					const [destItemRecord] = location.current.dropTargets;

					if (!isItem(destItemRecord.data)) return;

					const destItemId = destItemRecord.data.id;

					const indexOfTarget = itemsGetter().findIndex((item) => item.id === destItemId);

					const closestEdgeOfTarget = extractClosestEdge(destItemRecord.data);

					const destinationIndex = getReorderDestinationIndex({
						startIndex: itemIdx,
						indexOfTarget,
						closestEdgeOfTarget,
						axis: 'vertical'
					});

					reorder(itemId, itemIdx, destinationIndex);
				}
			});
		});
	}

	private handleKeydown(el: HTMLElement, e: KeyboardEvent) {
		if (e.key === 'Escape') {
			e.preventDefault();
			this.selectedItem = { id: '', index: -1 };

			this.mode = 'drag';
			el.blur();
			return;
		}

		if (this.mode === 'drag') {
			if (e.key === 'Enter') {
				this.mode = 'keyboard';
				this.selectedItem.index = 0;
			}
			return;
		}

		if (e.key === 'Enter') {
			this.selectedItem = { id: '', index: -1 };
			this.mode = 'drag';

			return;
		}

		if (e.key === ' ') {
			e.preventDefault();

			const targetId = this.#itemsGetter()[this.selectedItem.index].id;
			if (this.selectedItem.id === targetId) {
				this.selectedItem.id = '';
			} else {
				this.selectedItem.id = targetId;
			}

			return;
		}

		if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
			e.preventDefault();

			if (this.selectedItem.index <= 0) {
				return;
			}

			if (this.selectedItem.id) {
				this.#reorder(this.selectedItem.id, this.selectedItem.index, this.selectedItem.index - 1);
			}

			if (this.selectedItem.index > 0) {
				this.selectedItem.index--;
			}

			return;
		}

		if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
			e.preventDefault();

			if (this.selectedItem.index >= this.#itemsGetter().length - 1) {
				return;
			}

			if (this.selectedItem.id) {
				this.#reorder(this.selectedItem.id, this.selectedItem.index, this.selectedItem.index + 1);
			}

			if (this.selectedItem.index < this.#itemsGetter().length - 1) {
				this.selectedItem.index++;
			}

			return;
		}
	}

	private handleBlur() {
		this.selectedItem = { id: '', index: -1 };
		this.mode = 'drag';
	}

	private setClosestEdge(id: string, edge: Edge | null) {
		const dropIndicator = this.#dropIndicators.get(id);
		if (!dropIndicator) {
			return;
		}

		if (dropIndicator.edge === edge) {
			return;
		}

		dropIndicator.edge = edge;
		if (edge === null) {
			dropIndicator.el.removeAttribute('data-visible');
			dropIndicator.el.removeAttribute('data-closest-edge');
		} else {
			dropIndicator.el.setAttribute('data-visible', 'true');
			dropIndicator.el.setAttribute('data-closest-edge', edge);
		}
	}

	root(el: HTMLElement) {
		el.tabIndex = 0;
		el.setAttribute('data-part', 'root');

		const handleKeydown = (e: KeyboardEvent) => this.handleKeydown(el, e);

		$effect(() => {
			if (this.enabled) {
				el.setAttribute('data-enabled', 'true');
			} else {
				el.removeAttribute('data-enabled');
			}
		});

		$effect(() => {
			el.addEventListener('keydown', handleKeydown);
			el.addEventListener('blur', this.handleBlur);

			return () => {
				el.removeEventListener('keydown', handleKeydown);
				el.removeEventListener('blur', this.handleBlur);
			};
		});
	}

	item(el: HTMLElement, { id }: PartOptions) {
		el.setAttribute('data-part', 'item');

		$effect(() => {
			if (this.enabled) {
				el.setAttribute('data-enabled', 'true');
			} else {
				el.removeAttribute('data-enabled');
			}
		});

		$effect(() => {
			const cleanup = combine(
				draggable({
					element: el,
					dragHandle: this.#dragHandles.get(id),
					canDrag: () => this.enabled,
					getInitialData: () => {
						return getItem({ instanceId: this.#instanceId, id });
					},
					onDragStart() {
						el.setAttribute('data-dragging', 'true');
					},
					onDrop() {
						el.setAttribute('data-dragging', 'false');

						el.animate([{ backgroundColor: '#37cdbe' }, {}], {
							duration: 600,
							delay: 150,
							easing: 'cubic-bezier(0.25, 0.1, 0.25, 1.0)',
							iterations: 1
						});
					}
				}),
				dropTargetForElements({
					element: el,
					canDrop: ({ source }) => {
						return (
							isItem(source.data) && source.data.instanceId === this.#instanceId && this.enabled
						);
					},
					getData: ({ element, input }) => {
						return attachClosestEdge(getItem({ instanceId: this.#instanceId, id }), {
							element: element,
							input: input,
							allowedEdges: ['top', 'bottom']
						});
					},
					onDragEnter: ({ source, self }) => {
						if (isItem(source.data) && source.data.id !== id) {
							const edge = extractClosestEdge(self.data);
							this.setClosestEdge(id, edge);

							if (edge) {
								el.setAttribute('data-closest-edge', edge);
							} else {
								el.removeAttribute('data-closest-edge');
							}
						}
					},
					onDrag: ({ source, self }) => {
						if (isItem(source.data) && source.data.id !== id) {
							const edge = extractClosestEdge(self.data);
							this.setClosestEdge(id, edge);

							if (edge) {
								el.setAttribute('data-closest-edge', edge);
							} else {
								el.removeAttribute('data-closest-edge');
							}
						}
					},
					onDragLeave: () => {
						this.setClosestEdge(id, null);
						el.removeAttribute('data-closest-edge');
					},
					onDrop: () => {
						this.setClosestEdge(id, null);
						el.removeAttribute('data-closest-edge');
					}
				})
			);

			return () => {
				cleanup();
			};
		});
	}

	dragHandle(el: HTMLElement, { id }: PartOptions) {
		el.setAttribute('data-part', 'drag-handle');
		this.#dragHandles.set(id, el);

		$effect(() => {
			if (this.enabled) {
				el.setAttribute('data-enabled', 'true');
			} else {
				el.removeAttribute('data-enabled');
			}
		});
	}

	dropIndicator(el: HTMLElement, { id }: PartOptions) {
		el.setAttribute('data-part', 'drop-indicator');
		this.#dropIndicators.set(id, { el, edge: null });

		$effect(() => {
			if (this.enabled) {
				el.setAttribute('data-enabled', 'true');
			} else {
				el.removeAttribute('data-enabled');
			}
		});
	}
}
