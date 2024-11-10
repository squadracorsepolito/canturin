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

type Options = {
	instanceId: string;
	id: string;
	dragHandle?: HTMLElement;
	setClosestEdge: (edge: Edge | null) => void;
	setIsDragging?: (isDragging: boolean) => void;
};

export const sortableItem = (
	el: HTMLElement,
	{ instanceId, id, dragHandle, setClosestEdge, setIsDragging }: Options
) => {
	$effect(() => {
		const cleanup = combine(
			draggable({
				element: el,
				dragHandle: dragHandle,
				getInitialData() {
					return getItem({ instanceId, id });
				},
				onDragStart() {
					setIsDragging?.(true);
				},
				onDrop() {
					setIsDragging?.(false);

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
				canDrop({ source }) {
					return isItem(source.data) && source.data.instanceId === instanceId;
				},
				getData({ element, input }) {
					return attachClosestEdge(getItem({ instanceId, id }), {
						element: element,
						input: input,
						allowedEdges: ['top', 'bottom']
					});
				},
				onDragEnter({ source, self }) {
					if (isItem(source.data) && source.data.id !== id) {
						setClosestEdge(extractClosestEdge(self.data));
					}
				},
				onDrag({ source, self }) {
					if (isItem(source.data) && source.data.id !== id) {
						setClosestEdge(extractClosestEdge(self.data));
					}
				},
				onDragLeave() {
					setClosestEdge(null);
				},
				onDrop() {
					setClosestEdge(null);
				}
			})
		);

		return () => {
			cleanup();
		};
	});
};

export type Item = {
	[privateItemKey]: true;
	instanceId: string;
	id: string;
};

export function getItem(data: Omit<Item, typeof privateItemKey>) {
	return {
		[privateItemKey]: true,
		...data
	};
}

export function isItem(data: Record<string | symbol, unknown>): data is Item {
	return Boolean(data[privateItemKey]);
}

export function mountSortableWrapper<T extends { id: string }>(
	instanceId: string,
	itemsGetter: () => T[],
	reorder: (id: string, from: number, to: number) => void
) {
	onMount(() => {
		return monitorForElements({
			canMonitor({ source }) {
				return isItem(source.data) && source.data.instanceId === instanceId;
			},
			onDrop({ source, location }) {
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
