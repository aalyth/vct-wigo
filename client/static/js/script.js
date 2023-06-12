const url_bar = document.querySelector('#url') 
const depth_bar = document.querySelector('#depth')
const submit = document.querySelector('#submit')

const error = document.querySelector('.error')
const results = document.querySelector('.results-container')

submit.addEventListener('click', handleRequest)

url_bar.addEventListener('keypress', handleEnter)
depth_bar.addEventListener('keypress', handleEnter)

window.onload = () => {
	url_bar.value = ''
	depth_bar.value = ''
}

function handleEnter(event) {
	if (event.key === 'Enter') {
		handleRequest()
	}
}

async function handleRequest() {
	const url = (url_bar.value).replace(/\s/g, '_')
	const depth = parseInt(depth_bar.value)

	if (depth < 1 || depth > 3 || isNaN(depth)) { 
		error.innerHTML = 'Invalid depth value.'
		return
	}

	const resp = await fetch(`/api/wiki?url=wiki/${url}&depth=${depth}`);
	if (resp.status != 200) {
		error.innerHTML = 'Could not process request.'
		return
	}

	const pages = await resp.json();
	if (pages == null) {
		error.innerHTML = 'Invalid search.'
		return
	}
	error.innerHTML = ''
	results.innerHTML = ''
	console.log(pages)
	for (var page of pages) {
		var summary = page.Summary[0]
		if (page.Summary[1] != undefined) summary += page.Summary[1]
		results.insertAdjacentHTML('beforeend', `
			<div class="result">
				<a href="${page.Url}">${page.Title}</a>
				<p>${summary}</a>
			</div>
		`)
	}
	results.scrollTo(0, 0)
}

