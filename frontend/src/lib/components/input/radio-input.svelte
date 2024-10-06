<script lang="ts">
	import { TickIcon } from '../icon';
	import type { RadioInputOption } from './types';

	type Props = {
		name: string;
		selected: number;
		options: RadioInputOption[];
	};

	let { name, selected = $bindable(), options }: Props = $props();

	function getLabelClass(opt: RadioInputOption) {
		let cn = 'flex items-center gap-3 border-2 py-2 px-3 rounded-btn border-2 ';

		if (opt.disabled) {
			return cn + 'cursor-not-allowed bg-base-300 border-base-300 opacity-50';
		}

		if (opt.id === selected) {
			cn += 'text-primary border-primary ';
		} else {
			cn += 'border-base-300 ';
		}

		return cn + 'hover:cursor-pointer hover:bg-base-200';
	}
</script>

<ul class="flex flex-col gap-2">
	{#each options as opt}
		{@const inputId = name + '-' + opt.name}
		<li>
			<input
				type="radio"
				{name}
				bind:group={selected}
				id={inputId}
				disabled={opt.disabled}
				value={opt.id}
				class="hidden"
			/>

			<label for={inputId} class={getLabelClass(opt)}>
				<div class="flex-1">
					<div class="font-medium">{opt.label}</div>

					{#if opt.desc}
						<div class="text-xs opacity-90">{opt.desc}</div>
					{/if}
				</div>

				{#if opt.id === selected}
					<TickIcon />
				{/if}
			</label>
		</li>
	{/each}
</ul>
