<script lang="ts">
	import type { Message } from '$lib/api/canturin';
	import SignalGrid from '$lib/components/grid/signal-grid.svelte';
	import Table from '$lib/components/table/table.svelte';
	import { genSignalColor, getSignalKindName } from '$lib/utils';
	import Panel from './panel.svelte';

	type Props = {
		message: Message;
	};

	let { message }: Props = $props();
</script>

<Panel>
	<section>
		<h3>{message.name}</h3>
		<p>{message.desc}</p>
	</section>

	{#if message.signals}
		<section>
			<Table items={message.signals}>
				{#snippet header()}
					<th></th>
					<th>Name</th>
					<th>Start Pos</th>
					<th>Size</th>
				{/snippet}

				{#snippet row({ id, name, startPos, size, kind })}
					<td>
						<span class="block h-10 w-10 rounded-box" style:background-color={genSignalColor(name)}
						></span>
					</td>
					<td>
						<div class="flex flex-col gap-1">
							<div>
								<a class="link" href="#{id}">
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

			<SignalGrid signals={message.signals} sizeByte={message.sizeByte} />
		</section>

		{#each message.signals as signal}
			<section id={signal.id}>
				<h4># {signal.name}</h4>
				<p>{signal.desc}</p>
			</section>
		{/each}
	{/if}
</Panel>
