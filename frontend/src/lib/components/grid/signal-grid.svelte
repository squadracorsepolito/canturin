<script lang="ts">
	import type { Signal } from '$lib/api/canturin';
	import { genSignalColor } from '$lib/utils';
	import { BackwardIcon, ForwardIcon } from '../icon';

	const gridWidth = 8;

	type Props = {
		signals: Signal[];
		sizeByte: number;
	};

	let { signals, sizeByte }: Props = $props();

	type Item = {
		id: string;
		colStart: number;
		colEnd: number;
		rowStart: number;
		rowEnd: number;
		name: string;
		continues?: boolean;
		follows?: boolean;
	};

	let signalItems: Item[] = $state([]);
	let hoveredID = $state('');

	$effect(() => {
		const items: Item[] = [];

		for (const sig of signals) {
			const startPos = sig.startPos;
			const size = sig.size;
			const name = sig.name;
			const id = sig.id;

			const nRows = Math.ceil(((startPos % gridWidth) + size) / gridWidth);

			const startRow =
				startPos % gridWidth === 0
					? Math.ceil(startPos / gridWidth) + 1
					: Math.ceil(startPos / gridWidth);

			const startCol = (startPos % gridWidth) + 1;

			if (startPos % 8 === 0 && size % 8 === 0) {
				// spans exactly over one/multiple rows
				items.push({
					colStart: 1,
					colEnd: gridWidth + 1,
					rowStart: startRow,
					rowEnd: startRow + nRows,
					name,
					id
				});
				continue;
			}

			if (nRows === 1) {
				items.push({
					colStart: startCol,
					colEnd: startCol + size,
					rowStart: startRow,
					rowEnd: startRow,
					name,
					id
				});
				continue;
			}

			if (nRows === 2) {
				items.push(
					{
						colStart: startCol,
						colEnd: gridWidth + 1,
						rowStart: startRow,
						rowEnd: startRow,
						name,
						id,
						continues: true
					},
					{
						colStart: 1,
						colEnd: size - (gridWidth - (startCol - 1)) + 1,
						rowStart: startRow + nRows - 1,
						rowEnd: startRow + nRows - 1,
						name,
						id,
						follows: true
					}
				);
				continue;
			}

			items.push({
				colStart: startCol,
				colEnd: gridWidth + 1,
				rowStart: startRow,
				rowEnd: startRow,
				name,
				id,
				continues: true
			});

			items.push({
				colStart: 1,
				colEnd: gridWidth + 1,
				rowStart: startRow + 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				continues: true,
				follows: true
			});

			items.push({
				colStart: 1,
				colEnd: size - (startCol - 1) - gridWidth * (nRows - 2) + 1,
				rowStart: startRow + nRows - 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				follows: true
			});
		}

		// flip colStart/colEnd because 0 is on the right
		for (let i = 0; i < items.length; i++) {
			const cs = items[i].colStart;
			if (cs < 5) {
				items[i].colStart = cs + 2 * (5 - cs);
			} else if (cs) {
				items[i].colStart = cs - 2 * (cs - 5);
			}

			const ce = items[i].colEnd;
			if (cs < 5) {
				items[i].colEnd = ce + 2 * (5 - ce);
			} else if (cs) {
				items[i].colEnd = ce - 2 * (ce - 5);
			}
		}

		signalItems = items;
	});

	function getPlaceholder(idx: number) {
		let nRow = Math.ceil(idx / gridWidth);
		if (idx % gridWidth === 0) {
			nRow++;
		}
		return gridWidth - idx + 2 * gridWidth * (nRow - 1) - 1;
	}
</script>

{#snippet item({ colStart, colEnd, rowStart, rowEnd, id, name, continues, follows }: Item)}
	<div
		onpointerenter={() => (hoveredID = id)}
		onpointerleave={() => (hoveredID = '')}
		class="rounded-box text-xs cursor-pointer hover:opacity-50 transition-all {id === hoveredID &&
			'opacity-50'}"
		style:grid-column-start={colStart}
		style:grid-column-end={colEnd}
		style:grid-row-start={rowStart}
		style:grid-row-end={rowEnd}
		style:background-color={genSignalColor(name)}
	>
		<a href="#{id}" class="flex p-3 h-full w-full">
			{#if continues}
				<span class="self-end"><ForwardIcon /></span>
			{/if}

			<span class="flex-1 self-center text-center">
				{#if hoveredID !== id}
					{name}
				{/if}
			</span>

			{#if follows}
				<span class="self-start"> <BackwardIcon /></span>
			{/if}
		</a>
	</div>
{/snippet}

<div class="relative">
	<div class="absolute h-full w-full block rounded-box bg-base-200"></div>

	<div
		class="absolute w-full h-full grid gap-3 p-4"
		style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
		style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
	>
		{#each { length: gridWidth * sizeByte } as _, idx}
			<div class="flex items-center rounded-box justify-center bg-base-100">
				<span class="text-base-content font-medium">
					{getPlaceholder(idx)}
				</span>
			</div>
		{/each}
	</div>

	<div
		class="relative grid gap-3 p-4"
		style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
		style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
		style:aspect-ratio={gridWidth / sizeByte}
	>
		{#each signalItems as sigItem}
			{@render item(sigItem)}
		{/each}
	</div>
</div>
