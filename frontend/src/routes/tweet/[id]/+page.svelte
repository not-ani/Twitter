<script lang="ts">
	import { page } from '$app/stores';
	import type { InnerTweets } from '../../../types';

	$: id = $page.params.id;
	const fetchData = async () => {
		const response = await fetch(`http://localhost:3000/api/tweet/${id}/`);
		return (await response.json()) as InnerTweets;
	};
</script>

{#await fetchData()}
	<p>loading...</p>
{:then data}
	{data.content}
{:catch error}
	<p>error: {error.message}</p>
{/await}
