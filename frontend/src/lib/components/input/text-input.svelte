<script lang="ts">
	import { selectTextOnFocus } from '$lib/actions';

	let value = $state('initial value');
	let width = $state(0);
	let widthString = $derived(`${width + 2 + 1}px`);
</script>

<div class="relative">
	<input
		type="text"
		use:selectTextOnFocus
		bind:value
		class="input input-primary font-medium text-xl border-transparent focus:border-solid"
		style:width={widthString}
	/>

	<span
		class="px-4 inline-block absolute left-0 break-all top-12 font-medium text-xl"
		bind:clientWidth={width}
	>
		{value}
	</span>
</div>
