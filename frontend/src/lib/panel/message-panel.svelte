<script lang="ts">
	import { type Message } from '$lib/api/canturin';
	import { MessageByteOrder, SignalKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import { BackwardIcon, ForwardIcon } from '$lib/components/icon';
	import Table from '$lib/components/table/table.svelte';
	import randomColor from 'randomcolor';
	import Panel from './panel.svelte';
	import { getColorByName, getHexNumber } from '$lib/utils';
	import { useMessage, type GridItem } from '$lib/state/message-state.svelte';
	import type { SummaryInfo } from '$lib/components/summary/types';
	import Summary from '$lib/components/summary/summary.svelte';

	const gridWidth = 8;

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let state = useMessage(entityId);

	$effect(() => {
		state.reload(entityId);
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

	function getSummaryInfos(msg: Message) {
		const infos: SummaryInfo[] = [];

		infos.push({
			title: 'CAN ID',
			value: getHexNumber(msg.canId),
			desc: `Decimal: ${msg.canId}`
		});

		if (msg.hasStaticCANID) {
			infos[0].badge = { text: 'Static', color: 'secondary' };
		} else {
			infos[0].badge = { text: 'Generated', color: 'primary' };

			infos.push({
				title: 'Message ID',
				value: getHexNumber(msg.id),
				desc: `Decimal: ${msg.id}`
			});
		}

		infos.push({
			title: 'Size',
			value: msg.sizeByte,
			desc: `The size in bytes`
		});
		infos.push({
			title: 'Byte Order',
			value: getByteOrderName(msg.byteOrder),
			desc: `Endianness`
		});

		return infos;
	}
</script>

{#snippet gridItem({ colStart, colEnd, rowStart, rowEnd, id, name, continues, follows }: GridItem)}
	{@const color = getColorByName(name)}

	<div
		onpointerenter={() => (state.hoveredID = id)}
		onpointerleave={() => (state.hoveredID = '')}
		class="rounded-box text-xs cursor-pointer hover:opacity-50 transition-all {id ===
			state.hoveredID && 'opacity-50'}"
		style:grid-column-start={colStart}
		style:grid-column-end={colEnd}
		style:grid-row-start={rowStart}
		style:grid-row-end={rowEnd}
		style:background-color={color.bgColor}
		style:color={color.textColor}
	>
		<div class="flex p-3 h-full w-full">
			{#if state.hoveredID !== id}
				{#if continues}
					<span class="self-end"><ForwardIcon /></span>
				{/if}

				<span class="flex-1 self-center text-center truncate text-ellipsis">
					{name}
				</span>

				{#if follows}
					<span class="self-start"> <BackwardIcon /></span>
				{/if}
			{/if}
		</div>
	</div>
{/snippet}

{#snippet grid(sizeByte: number, items: GridItem[])}
	<div class="flex gap-2">
		<div class="flex flex-col py-2 gap-2 justify-around">
			{#each { length: sizeByte } as _, idx}
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
				style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
			>
				{#each { length: gridWidth * sizeByte } as _, idx}
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
				style:grid-template-rows="repeat({sizeByte}, minmax(0, 1fr))"
				style:aspect-ratio={gridWidth / sizeByte}
			>
				{#each items as sigItem}
					{@render gridItem(sigItem)}
				{/each}
			</div>
		</div>
	</div>
{/snippet}

{#snippet msgPanel(msg: Message)}
	<section>
		<h3>{msg.name}</h3>
		<p>{msg.desc}</p>
	</section>

	<section>
		<Summary infos={getSummaryInfos(msg)} />
	</section>

	{#if msg.signals}
		<section class="flex flex-col gap-2 @5xl:flex-row">
			<div class="flex-1">
				<Table items={msg.signals}>
					{#snippet header()}
						<th></th>
						<th>Name</th>
						<th>Start Pos</th>
						<th>Size</th>
					{/snippet}

					{#snippet row({ entityId, name, startPos, size, kind })}
						<td>
							<span
								class="block h-10 w-10 rounded-box {state.hoveredID === entityId &&
									'scale-125'} transition-all"
								style:background-color={genSignalColor(name)}
							></span>
						</td>
						<td>
							<div class="flex flex-col gap-1">
								<div
									onpointerenter={() => (state.hoveredID = entityId)}
									onpointerleave={() => (state.hoveredID = '')}
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
				{@render grid(msg.sizeByte, state.gridItems)}
			</div>
		</section>

		{#each msg.signals as signal}
			<section id={signal.entityId}>
				<h4># {signal.name}</h4>
				<p>{signal.desc}</p>
			</section>
		{/each}
	{/if}
{/snippet}

<Panel>
	{#if !state.isLoading && state.message}
		{@render msgPanel(state.message)}
	{/if}
</Panel>
