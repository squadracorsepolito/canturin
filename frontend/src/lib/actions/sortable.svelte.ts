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

enum Keys {
	enter = 'Enter',
	escape = 'Escape',
	space = ' ',
	arrowLeft = 'ArrowLeft',
	arrowRight = 'ArrowRight',
	arrowUp = 'ArrowUp',
	arrowDown = 'ArrowDown'
}

enum ItemStates {
	idle = 'idle',
	dragging = 'dragging',
	selecting = 'selecting',
	moving = 'moving'
}

enum Mode {
	drag = 'drag',
	keyboard = 'keyboard'
}

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

	#items = new Map<string, HTMLElement>();
	#dragHandles = new Map<string, HTMLElement>();
	#dropIndicator: HTMLElement | undefined = undefined;

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

	private handleMode(el: HTMLElement, mode: Mode) {
		this.#mode = mode;
		el.setAttribute('data-mode', mode);
	}

	private handleItemState(el: HTMLElement, state: ItemStates) {
		el.setAttribute('data-state', state);
	}

	private resetItemStates() {
		this.#selectedItemIdx = -1;
		this.#movingItemId = '';

		for (const item of this.#items.values()) {
			this.handleItemState(item, ItemStates.idle);
		}
	}

	#mode: Mode = Mode.drag;
	#selectedItemIdx = -1;
	#movingItemId = $state('');

	private handleKeydown(el: HTMLElement, e: KeyboardEvent) {
		if (e.key === Keys.escape) {
			e.preventDefault();
			this.resetItemStates();
			this.handleMode(el, Mode.drag);

			this.enabled = false;

			el.blur();

			return;
		}

		if (!this.enabled) {
			return;
		}

		const items = this.#itemsGetter();
		if (items.length === 0) return;

		if (e.key === Keys.enter) {
			if (this.#mode === Mode.drag) {
				const firstItemId = items[0].id;
				const firstItem = this.#items.get(firstItemId);
				if (!firstItem) return;

				this.handleMode(el, Mode.keyboard);
				this.handleItemState(firstItem, ItemStates.selecting);
				this.#selectedItemIdx = 0;
			} else {
				this.resetItemStates();
				this.handleMode(el, Mode.drag);
			}

			return;
		}

		if (e.key === Keys.space) {
			e.preventDefault();

			const itemId = items[this.#selectedItemIdx].id;
			const item = this.#items.get(itemId);
			if (!item) return;

			if (this.#movingItemId === itemId) {
				this.handleItemState(item, ItemStates.selecting);
				this.#movingItemId = '';
			} else {
				this.handleItemState(item, ItemStates.moving);
				this.#movingItemId = itemId;
			}

			return;
		}

		if (e.key === Keys.arrowUp || e.key === Keys.arrowLeft) {
			e.preventDefault();

			// reached the first item
			if (this.#selectedItemIdx <= 0) {
				return;
			}

			// changing the selecting item
			if (!this.#movingItemId) {
				const prevItem = this.#items.get(items[this.#selectedItemIdx].id);
				if (!prevItem) return;

				const nextItem = this.#items.get(items[this.#selectedItemIdx - 1].id);
				if (!nextItem) return;

				this.handleItemState(prevItem, ItemStates.idle);
				this.handleItemState(nextItem, ItemStates.selecting);

				this.#selectedItemIdx--;

				return;
			}

			// swapping the moving item
			const itemId = items[this.#selectedItemIdx].id;
			const item = this.#items.get(itemId);
			if (!item) return;

			const fromIdx = this.#selectedItemIdx;
			const toIdx = this.#selectedItemIdx - 1;
			this.#reorder(itemId, fromIdx, toIdx);

			this.#selectedItemIdx--;

			return;
		}

		if (e.key === Keys.arrowDown || e.key === Keys.arrowRight) {
			e.preventDefault();

			// reached the last item
			if (this.#selectedItemIdx >= items.length - 1) {
				return;
			}

			// changing the selecting item
			if (!this.#movingItemId) {
				const prevItem = this.#items.get(items[this.#selectedItemIdx].id);
				if (!prevItem) return;

				const nextItem = this.#items.get(items[this.#selectedItemIdx + 1].id);
				if (!nextItem) return;

				this.handleItemState(prevItem, ItemStates.idle);
				this.handleItemState(nextItem, ItemStates.selecting);

				this.#selectedItemIdx++;

				return;
			}

			// swapping the moving item
			const itemId = items[this.#selectedItemIdx].id;
			const item = this.#items.get(itemId);
			if (!item) return;

			const fromIdx = this.#selectedItemIdx;
			const toIdx = this.#selectedItemIdx + 1;
			this.#reorder(itemId, fromIdx, toIdx);

			this.#selectedItemIdx++;

			return;
		}
	}

	private handleClosestEdge(itemEl: HTMLElement, edge: Edge | null) {
		if (edge === null) {
			itemEl.removeAttribute('data-closest-edge');
		} else {
			itemEl.setAttribute('data-closest-edge', edge);
		}

		if (!this.#dropIndicator) {
			return;
		}

		if (edge === null) {
			this.#dropIndicator.hidden = true;
			this.#dropIndicator.removeAttribute('data-visible');
			return;
		}

		this.#dropIndicator.hidden = false;
		this.#dropIndicator.setAttribute('data-visible', 'true');

		const itemRect = itemEl.getBoundingClientRect();
		const indicatorRect = this.#dropIndicator.getBoundingClientRect();

		this.#dropIndicator.style.left = `${itemRect.left}px`;

		if (edge === 'top') {
			this.#dropIndicator.style.top = `${itemRect.top - indicatorRect.height / 2}px`;
		} else {
			this.#dropIndicator.style.top = `${itemRect.top + itemRect.height - indicatorRect.height / 2}px`;
		}
	}

	private handleEnabled(el: HTMLElement, enabled: boolean) {
		if (enabled) {
			el.setAttribute('data-enabled', 'true');
		} else {
			el.removeAttribute('data-enabled');
		}
	}

	private initPart(el: HTMLElement, partName: string) {
		el.setAttribute('data-scope', 'sortable');
		el.setAttribute('data-part', partName);
	}

	root(el: HTMLElement) {
		this.initPart(el, 'root');

		el.tabIndex = 0;

		const handleKeydown = (e: KeyboardEvent) => this.handleKeydown(el, e);

		const handleBlur = () => {
			this.resetItemStates();
		};

		$effect(() => {
			if (this.enabled) {
				el.setAttribute('data-enabled', 'true');
				el.focus();
			} else {
				el.removeAttribute('data-enabled');
			}
		});

		$effect(() => {
			el.addEventListener('keydown', handleKeydown);
			el.addEventListener('blur', handleBlur);

			return () => {
				el.removeEventListener('keydown', handleKeydown);
				el.removeEventListener('blur', handleBlur);
			};
		});
	}

	item(el: HTMLElement, { id }: PartOptions) {
		this.initPart(el, 'item');
		this.handleItemState(el, ItemStates.idle);

		this.#items.set(id, el);

		let prevEdge: Edge | null = null;

		$effect(() => {
			this.handleEnabled(el, this.enabled);
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
					onDragStart: () => {
						this.handleItemState(el, ItemStates.dragging);
					},
					onDrop: () => {
						this.handleItemState(el, ItemStates.idle);

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
							if (edge === prevEdge) {
								return;
							}

							prevEdge = edge;
							this.handleClosestEdge(el, edge);
						}
					},
					onDrag: ({ source, self }) => {
						if (isItem(source.data) && source.data.id !== id) {
							const edge = extractClosestEdge(self.data);
							if (edge === prevEdge) {
								return;
							}

							prevEdge = edge;
							this.handleClosestEdge(el, edge);
						}
					},
					onDragLeave: () => {
						prevEdge = null;
						this.handleClosestEdge(el, null);
					},
					onDrop: () => {
						prevEdge = null;
						this.handleClosestEdge(el, null);
					}
				})
			);

			return () => {
				cleanup();
				this.#items.delete(id);
			};
		});
	}

	dragHandle(el: HTMLElement, { id }: PartOptions) {
		this.initPart(el, 'drag-handle');

		this.#dragHandles.set(id, el);

		$effect(() => {
			this.handleEnabled(el, this.enabled);
		});
	}

	dropIndicator(el: HTMLElement) {
		this.initPart(el, 'drop-indicator');

		el.style.position = 'absolute';
		el.hidden = true;

		this.#dropIndicator = el;

		$effect(() => {
			this.handleEnabled(el, this.enabled);
		});
	}

	isItemMoving(id: string) {
		return this.#movingItemId === id;
	}
}
