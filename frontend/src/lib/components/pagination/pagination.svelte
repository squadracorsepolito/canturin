<script lang="ts">
	import { AltArrowLeftIcon, AltArrowRightIcon } from '../icon';
	import * as numberInput from '@zag-js/number-input';
	import { normalizeProps, useMachine } from '@zag-js/svelte';

	type Props = {
		selectedPage: number;
		totPages: number;
		pagesPerRow?: number;
	};

	let { selectedPage = $bindable(), totPages, pagesPerRow = 8 }: Props = $props();

	function getSectionFromPage(page: number) {
		return Math.floor(page / pagesPerRow);
	}

	let currSection = $state(getSectionFromPage(selectedPage));

	const id = $props.id();
	const service = useMachine(numberInput.machine, () => ({
		id,
		defaultValue: selectedPage.toString(),
		min: 0,
		max: totPages - 1,
		clampValueOnBlur: true,
		onFocusChange: (details) => {
			if (details.focused) return;

			let tmpVal = details.valueAsNumber;
			if (isNaN(tmpVal) || tmpVal < 0) {
				tmpVal = 0;
			} else if (tmpVal >= totPages) {
				tmpVal = totPages - 1;
			}

			currSection = getSectionFromPage(tmpVal);

			selectedPage = tmpVal;
		}
	}));

	const api = $derived(numberInput.connect(service, normalizeProps));

	function handleSectionChange(offset: number) {
		currSection = currSection + offset;

		selectedPage = currSection * pagesPerRow;
		if (offset < 0) {
			selectedPage += pagesPerRow - 1;
		}

		if (selectedPage >= totPages) {
			selectedPage = totPages - 1;
		}

		api.setValue(selectedPage);
	}

	function handlePageClick(page: number) {
		selectedPage = page;
		api.setValue(page);
	}
</script>

{#snippet input()}
	<form class="join-item bg-base-200 flex items-center">
		<div {...api.getRootProps()}>
			<input {...api.getInputProps()} />
		</div>
	</form>
{/snippet}

<div class="join">
	<button
		onclick={() => handleSectionChange(-1)}
		disabled={currSection === 0}
		class="join-item btn"
	>
		<AltArrowLeftIcon />
	</button>

	{#each { length: pagesPerRow } as _, idx}
		{#if idx === pagesPerRow / 2}
			{@render input()}
		{/if}

		{@const currIdx = currSection * pagesPerRow + idx}
		<button
			onclick={() => handlePageClick(currIdx)}
			class="join-item btn {selectedPage === currIdx && 'btn-primary'}"
			disabled={currIdx >= totPages}>{currIdx}</button
		>
	{/each}

	<button
		onclick={() => handleSectionChange(1)}
		disabled={currSection === getSectionFromPage(totPages)}
		class="join-item btn"
	>
		<AltArrowRightIcon />
	</button>
</div>

<style lang="postcss">
	[data-scope='number-input'][data-part='root'] {
		@apply border-2 mx-2 rounded-btn border-transparent transition-colors;

		&[data-focus] {
			@apply border-primary focus-ring-primary;
		}
	}

	[data-scope='number-input'][data-part='input'] {
		@apply p-1 max-w-24 text-center outline-none rounded-btn;
	}
</style>
