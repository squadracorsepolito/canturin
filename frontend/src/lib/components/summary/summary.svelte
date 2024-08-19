<script lang="ts">
	import type { SummaryInfo } from './types';

	type Props = {
		infos: SummaryInfo[];
	};

	let { infos }: Props = $props();

	function getBadgeColor(color: 'primary' | 'secondary') {
		switch (color) {
			case 'primary':
				return 'badge-primary';
			case 'secondary':
				return 'badge-secondary';
		}
	}
</script>

<div class="stats">
	{#each infos as info}
		<div class="stat place-items-center">
			<div class="stat-title">
				{#if info.badge}
					<span class="badge {getBadgeColor(info.badge.color)}">
						{info.badge.text}
					</span>

					<span>{info.title}</span>
				{:else}
					{info.title}
				{/if}
			</div>

			<div class="stat-value">{info.value}</div>

			{#if info.desc}
				<div class="stat-desc">{info.desc}</div>
			{/if}
		</div>
	{/each}
</div>
