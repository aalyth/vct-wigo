const url_bar = document.querySelector('#url') 
const depth_bar = document.querySelector('#depth')
const submit = document.querySelector('#submit')

submit.addEventListener('click', handleRequest)

url_bar.addEventListener('keypress', handleEnter)
depth_bar.addEventListener('keypress', handleEnter)

function handleEnter(event) {
	if (event.key === 'Enter') {
		handleRequest()
	}
}

async function handleRequest() {
	const url = url_bar.value
	const depth = int(depth_bar.value)
	console.log('url = ', url)
	console.log('depth = ', depth)

	if (depth < 1 || depth > 3) { 
		alert('Error: invalid depth value.')
		return
	}

	const resp = await fetch(`/api/wiki?url=${url}&depth=${depth}`);
	const res = await resp.json();
	console.log(res);
}

