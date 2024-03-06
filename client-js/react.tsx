import React from 'react';
import r2wc from "react-to-webcomponent"
import ReactDOM from "react-dom/client";

console.log("Initializing React");

type Props = {
    cls?: string;
};

const TestReactComponent = ({ cls = "" }: Props) => {
    return <div className={`p-4 bg-red-50 rounded border-2 border-red-300 ${cls}`}>Test React Component</div>;
};

const TestReactComponentWC = r2wc(TestReactComponent, React, ReactDOM, {
    props: {
        cls: "string",
    },
});

customElements.define("test-rc", TestReactComponentWC);

