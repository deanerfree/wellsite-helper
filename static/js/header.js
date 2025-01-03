class HeaderComponent extends HTMLElement {
  constructor () {
    super()
    // this.attachShadow({ mode: 'open' })
    this.innerHTML = `
      <link rel="stylesheet" href="/static/styles.css">
      <style>
        header {
          background-color: #333;
          color: white;
          padding: 1rem;
          text-align: center;
          width: 100%;
        }
      </style>
      <header>
        <h1 class="text-4xl text-white uppercase font-semibold">Wellsite helper</h1>
        <nav class="flex justify-center space-x-4">
          <a href="/" class="text-blue-500">Home</a>
          <a href="/wellpath" class="text-blue-500">Wellpath</a>
          <a href="/about" class="text-blue-500">About</a>
          <a href="/contact" class="text-blue-500">Contact</a>
        </nav>
      </header>
    `
  }
}

customElements.define('header-component', HeaderComponent)
