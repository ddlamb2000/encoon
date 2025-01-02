<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<div class="text-column">
		<h1>εncooη</h1>
		<p>For the time being, εncooη is an experimental web-based application 
		oriented toward data management.</p>
		<p>The object of the reasearches εncooη is based on is to propose 
		a true engine that allows data structuration, data relationship 
		management, data presentation, and, mainly, to make obvious data navigation.</p>	
		<p>This brand new application doesn't require database system 
		learning or skills, because everything in the application 
		is managed through a very simple web-based user interface, 
		in real time. </p>
		<p>Primary dedicated to small business structures, 
		the application lets manage a large amount of multi-purposes information 
		within the organization, and share this information across business stakeholders 
		in an easy and natural way, using one Internet browser only.</p>
		<p>εncooη is also the foundation for the development of a true business-oriented 
		application development toolkit.</p>
		<p>Copyright David Lambert 2024</p>
	</div>
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}

	.welcome {
		display: block;
		position: relative;
		width: 100%;
		height: 0;
		padding: 0 0 calc(100% * 495 / 2048) 0;
	}

	.welcome img {
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		display: block;
	}
</style>
