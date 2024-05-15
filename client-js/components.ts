declare const up: any;

console.log("Initializing components");

type AlpineComponentData = {
    open: boolean;
};

up.compiler('.comp__alpine_test', (el: Element, data: AlpineComponentData) => {
    const button = el.querySelector('button');
    const view = el.querySelector('div');
    if (!view) return console.error('No view found');
    if (!button) return console.error('No button found');

    button.addEventListener('click', () => {
        data.open = !data.open;
        if (data.open) {
            button.classList.remove('bg-red-500');
            button.classList.add('bg-black');
            view.style.opacity = '0';
            view.classList.remove('hidden');
            up.animate(view, 'fade-in', {
                duration: 250,
                easing: 'linear'
              });
            view.style.removeProperty('opacity');
        } else {
            button.classList.add('bg-red-500');
            button.classList.remove('bg-black');
            view.classList.add('hidden');
        }
    });
});
