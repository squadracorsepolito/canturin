<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { Signal } from '$lib/api/canturin';
	import { BackwardIcon, ForwardIcon } from '../icon';

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

	type Props = {
		signals: Signal[];
		height: number;
		width?: number;
		onclick?: (id: string) => void;
	};

	let { signals, height, width = 8, onclick }: Props = $props();

	let hoveredId = $state('');

	let items = $derived.by(() => {
		if (signals.length === 0) {
			return [];
		}

		const res: Item[] = [];

		for (const sig of signals) {
			const startPos = sig.startPos;
			const size = sig.size;
			const name = sig.name;
			const id = sig.entityId;

			const nRows = Math.ceil(((startPos % width) + size) / width);

			const startRow =
				startPos % width === 0 ? Math.ceil(startPos / width) + 1 : Math.ceil(startPos / width);

			const startCol = (startPos % width) + 1;

			if (startPos % 8 === 0 && size % 8 === 0) {
				// spans exactly over one/multiple rows
				res.push({
					colStart: 1,
					colEnd: width + 1,
					rowStart: startRow,
					rowEnd: startRow + nRows,
					name,
					id
				});
				continue;
			}

			if (nRows === 1) {
				res.push({
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
				res.push(
					{
						colStart: startCol,
						colEnd: width + 1,
						rowStart: startRow,
						rowEnd: startRow,
						name,
						id,
						continues: true
					},
					{
						colStart: 1,
						colEnd: size - (width - (startCol - 1)) + 1,
						rowStart: startRow + nRows - 1,
						rowEnd: startRow + nRows - 1,
						name,
						id,
						follows: true
					}
				);
				continue;
			}

			res.push({
				colStart: startCol,
				colEnd: width + 1,
				rowStart: startRow,
				rowEnd: startRow,
				name,
				id,
				continues: true
			});

			res.push({
				colStart: 1,
				colEnd: width + 1,
				rowStart: startRow + 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				continues: true,
				follows: true
			});

			res.push({
				colStart: 1,
				colEnd: size - (startCol - 1) - width * (nRows - 2) + 1,
				rowStart: startRow + nRows - 1,
				rowEnd: startRow + nRows - 1,
				name,
				id,
				follows: true
			});
		}

		// flip colStart/colEnd because 0 is on the right
		for (let i = 0; i < res.length; i++) {
			const cs = res[i].colStart;
			if (cs < 5) {
				res[i].colStart = cs + 2 * (5 - cs);
			} else if (cs) {
				res[i].colStart = cs - 2 * (cs - 5);
			}

			const ce = res[i].colEnd;
			if (cs < 5) {
				res[i].colEnd = ce + 2 * (5 - ce);
			} else if (cs) {
				res[i].colEnd = ce - 2 * (ce - 5);
			}
		}

		return res;
	});

	function getPlaceholder(idx: number) {
		let nRow = Math.ceil(idx / width);
		if (idx % width === 0) {
			nRow++;
		}
		return width - idx + 2 * width * (nRow - 1) - 1;
	}
</script>

{#snippet item({ colStart, colEnd, id, rowStart, rowEnd, name, continues, follows }: Item)}
	<button
		use:colorByName={{ name }}
		onclick={() => onclick?.(id)}
		onpointerenter={() => (hoveredId = id)}
		onpointerleave={() => (hoveredId = '')}
		class="rounded-btn text-xs cursor-pointer hover:opacity-50 transition-all border-2"
		class:opacity-50={hoveredId === id}
		style:grid-column-start={colStart}
		style:grid-column-end={colEnd}
		style:grid-row-start={rowStart}
		style:grid-row-end={rowEnd}
	>
		<div class="flex p-2 h-full w-full">
			{#if hoveredId !== id}
				{#if continues}
					<span class="self-end"><ForwardIcon /></span>
				{/if}

				<span class="flex-1 self-center text-center truncate text-ellipsis">
					{name}
				</span>

				{#if follows}
					<span class="self-start"><BackwardIcon /></span>
				{/if}
			{/if}
		</div>
	</button>
{/snippet}

<div class="flex gap-2">
	<div class="flex flex-col gap-2 justify-around">
		{#each { length: height } as _, idx}
			<div class="font-medium pr-5">
				{idx + 1}
			</div>
		{/each}
	</div>

	<div class="flex-1 relative">
		<div class="absolute h-full w-full block rounded-box"></div>

		<div
			class="absolute w-full h-full grid gap-2"
			style:grid-template-columns="repeat({width}, minmax(0, 1fr))"
			style:grid-template-rows="repeat({height}, minmax(0, 1fr))"
		>
			{#each { length: width * height } as _, idx}
				<div
					class="flex items-center rounded-btn justify-center bg-base-200 border-2 border-base-300"
				>
					<span class="text-base-content font-medium">
						{getPlaceholder(idx)}
					</span>
				</div>
			{/each}
		</div>

		<div
			class="relative grid gap-2"
			style:grid-template-columns="repeat({width}, minmax(0, 1fr))"
			style:grid-template-rows="repeat({height}, minmax(0, 1fr))"
			style:aspect-ratio={width / height}
		>
			{#each items as sigItem (sigItem.id + sigItem.rowStart)}
				{@render item(sigItem)}
			{/each}
		</div>
	</div>
</div>
