<script lang="ts">
	import { MessageService, type Message } from '$lib/api/canturin';
	import { MessageByteOrder, SignalKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import { BackwardIcon, ForwardIcon } from '$lib/components/icon';
	import Table from '$lib/components/table/table.svelte';
	import randomColor from 'randomcolor';
	import Panel from './panel.svelte';
	import { getHexNumber } from '$lib/utils';
	import { Textarea } from '$lib/components/textarea';
	import { Editable } from '$lib/components/editable';

	const gridWidth = 8;

	type Props = {
		message: Message;
	};

	let { message }: Props = $props();

	type GridItem = {
		id: string;
		colStart: number;
		colEnd: number;
		rowStart: number;
		rowEnd: number;
		name: string;
		continues?: boolean;
		follows?: boolean;
	};

	let gridItems: GridItem[] = $state([]);
	let hoveredID = $state('');

	$effect(() => {
		if (!message.signals) {
			return;
		}

		const items: GridItem[] = [];

		for (const sig of message.signals) {
			const startPos = sig.startPos;
			const size = sig.size;
			const name = sig.name;
			const id = sig.entityId;

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

		gridItems = items;
	});

	function getPlaceholder(idx: number) {
		let nRow = Math.ceil(idx / gridWidth);
		if (idx % gridWidth === 0) {
			nRow++;
		}
		return gridWidth - idx + 2 * gridWidth * (nRow - 1) - 1;
	}

	function genSignalColor(signalName: string) {
		return randomColor({
			seed: signalName
		});
	}

	function getSignalKindName(sigKind: SignalKind): string {
		switch (sigKind) {
			case SignalKind.SignalKindStandard:
				return 'Standard';
			case SignalKind.SignalKindEnum:
				return 'Enum';
			case SignalKind.SignalKindMultiplexer:
				return 'Multiplexer';
		}
	}

	function getByteOrderName(byteOrder: MessageByteOrder): string {
		switch (byteOrder) {
			case MessageByteOrder.MessageByteOrderLittleEndian:
				return 'Little';
			case MessageByteOrder.MessageByteOrderBigEndian:
				return 'Big';
		}
	}

	let desc = $state('');
	$effect(() => {
		desc = message.desc;
	});

	let name = $state(message.name);
	$effect(() => {
		name = message.name;
	});

	async function setDesc() {
		try {
			const msg = await MessageService.SetDesc(desc);
			message = msg;
		} catch (error) {
			console.error(error);
		}
	}

	async function updateName(name: string) {
		try {
			const msg = await MessageService.UpdateName(name);
			message = msg;
		} catch (error) {
			console.error(error);
		}
	}
</script>

{#snippet gridItem({ colStart, colEnd, rowStart, rowEnd, id, name, continues, follows }: GridItem)}
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
		<a href="#{id}" class="link flex p-3 h-full w-full">
			{#if continues}
				<span class="self-end"><ForwardIcon /></span>
			{/if}

			<span class="flex-1 self-center text-center truncate text-ellipsis">
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

{#snippet grid()}
	<div class="flex gap-2">
		<div class="flex flex-col py-2 gap-2 justify-around">
			{#each { length: message.sizeByte } as _, idx}
				<div class="font-semibold pr-5">
					{idx + 1}
				</div>
			{/each}
		</div>

		<div class="flex-1 relative">
			<div class="absolute h-full w-full block rounded-box bg-base-200"></div>

			<div
				class="absolute w-full h-full grid gap-2 p-2"
				style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
				style:grid-template-rows="repeat({message.sizeByte}, minmax(0, 1fr))"
			>
				{#each { length: gridWidth * message.sizeByte } as _, idx}
					<div class="flex items-center rounded-box justify-center bg-base-100">
						<span class="text-base-content font-semibold">
							{getPlaceholder(idx)}
						</span>
					</div>
				{/each}
			</div>

			<div
				class="relative grid gap-2 p-2"
				style:grid-template-columns="repeat({gridWidth}, minmax(0, 1fr))"
				style:grid-template-rows="repeat({message.sizeByte}, minmax(0, 1fr))"
				style:aspect-ratio={gridWidth / message.sizeByte}
			>
				{#each gridItems as sigItem}
					{@render gridItem(sigItem)}
				{/each}
			</div>
		</div>
	</div>
{/snippet}

<Panel>
	<section>
		<h3>{message.name}</h3>
		<p>{message.desc}</p>
	</section>

	<Editable initialValue={message.name} onsubmit={(v) => updateName(v)}></Editable>

	<Textarea bind:value={desc}></Textarea>
	<button onclick={setDesc}>save desc</button>

	<section class="pb-5">
		<div class="stats">
			<div class="stat">
				<div class="stat-title">
					{#if message.hasStaticCANID}
						<span class="badge badge-secondary">Static</span>
					{:else}
						<span class="badge badge-primary">Generated</span>
					{/if}
					<span>CAN ID</span>
				</div>
				<div class="stat-value">{getHexNumber(message.canId)}</div>
				<div class="stat-desc">Decimal: {message.canId}</div>
			</div>

			{#if !message.hasStaticCANID}
				<div class="stat">
					<div class="stat-title">Message ID</div>
					<div class="stat-value">{getHexNumber(message.id)}</div>
					<div class="stat-desc">Decimal: {message.id}</div>
				</div>
			{/if}

			<div class="stat">
				<div class="stat-title">Size</div>
				<div class="stat-value">{message.sizeByte}</div>
				<div class="stat-desc">Bytes</div>
			</div>

			<div class="stat">
				<div class="stat-title">Byte Order</div>
				<div class="stat-value">{getByteOrderName(message.byteOrder)}</div>
				<div class="stat-desc">Endian</div>
			</div>
		</div>
	</section>

	{#if message.signals}
		<section class="flex flex-col gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={message.signals}>
					{#snippet header()}
						<th></th>
						<th>Name</th>
						<th>Start Pos</th>
						<th>Size</th>
					{/snippet}

					{#snippet row({ entityId, name, startPos, size, kind })}
						<td>
							<span
								class="block h-10 w-10 rounded-box {hoveredID === entityId &&
									'scale-125'} transition-all"
								style:background-color={genSignalColor(name)}
							></span>
						</td>
						<td>
							<div class="flex flex-col gap-1">
								<div
									onpointerenter={() => (hoveredID = entityId)}
									onpointerleave={() => (hoveredID = '')}
								>
									<a class="link" href="#{entityId}">
										<span>{name}</span>
									</a>
								</div>
								<span class="badge badge-ghost badge-sm">{getSignalKindName(kind)}</span>
							</div>
						</td>
						<td>{startPos}</td>
						<td>{size}</td>
					{/snippet}
				</Table>
			</div>

			<div class="@5xl:divider @5xl:divider-horizontal"></div>

			<div class="flex-1">
				{@render grid()}
			</div>
		</section>

		{#each message.signals as signal}
			<section id={signal.entityId}>
				<h4># {signal.name}</h4>
				<p>{signal.desc}</p>
			</section>
		{/each}
	{/if}
</Panel>
