const jsonTypes = ["string", "number", "boolean", "object", "array", "null"];
let customTypes = [];

function toPascalCase(str) {
    return str.replace(/(^\w|[a-zA-Z])\S*/g, match => {
        return match.charAt(0).toUpperCase() + match.slice(1).toLowerCase();
    }).replace(/\s+/g, '');
}

function addProperty(button, propName = '', propType = '', isArray = false) {
    const typeSection = button.parentElement;
    const propertySection = document.createElement("div");
    propertySection.className = "property-section";
    
    props = jsonTypes.concat(customTypes)
    //dirty hack, please revise
    if (!props.includes(propType)) {
        props.push(propType)
        console.log(`Hacked '${propType}' in...`)
    } else {
        console.log(`'${propType}' in there!`)
    }

    propType = propType.replace(/([A-Z])/g, ' $1').trim()
    props = props.map(t => t.replace(/([A-Z])/g, ' $1').trim())

    propertySection.innerHTML = `
        <label>Property name: <input type="text" class="property-name" value="${propName}"></label>
        <label>Type: 
            <select class="property-type">
                ${props.map(type => `<option value="${type}"${type === propType ? ' selected' : ''}>${type}</option>`).join('')}
            </select>
        </label>
        <label>Array: <input type="checkbox" class="property-array"${isArray ? ' checked' : ''}></label>
        <button onclick="removeProperty(this)">Remove property</button>
    `;


    const a = props.map(type => `${propType}|${type} ${type === propType}`)
    // console.log(`${props.map(type => `<option value="${type}"${type === propType ? ' selected' : ''}>${type}</option>`).join('')}`)
    console.log(a)

    typeSection.appendChild(propertySection);
}

function addType(typeName = '') {
    typeName = typeName || document.getElementById("new-type-name").value.trim();
    if (typeName && !customTypes.includes(typeName)) {
        customTypes.push(typeName);
        const typeSection = document.createElement("div");
        typeSection.className = "type-section";
        typeSection.setAttribute("data-type", typeName);
        typeSection.innerHTML = `
            <h3>${typeName} <button onclick="removeType(this)">Remove type</button></h3>
            <button onclick="addProperty(this)">Add property</button>
        `;
        document.getElementById("type-definitions").appendChild(typeSection);
        // updateTypeOptions();
        document.getElementById("new-type-name").value = '';
    }
}

function removeProperty(button) {
    const propertySection = button.parentElement;
    propertySection.parentElement.removeChild(propertySection);
}

function removeType(button) {
    const typeSection = button.parentElement.parentElement;
    const typeName = typeSection.getAttribute("data-type");
    customTypes.splice(customTypes.indexOf(typeName), 1);
    typeSection.parentElement.removeChild(typeSection);
    updateTypeOptions();
}


function updateTypeOptions() {
    document.querySelectorAll(".property-type").forEach(select => {
        const selectedValue = select.value;
        console.log(`Previous selected ${select.value}`)
        select.innerHTML = jsonTypes.concat(customTypes).map(type => `<option value="${type}">${type}</option>`).join('');
        select.value = selectedValue;
        // console.log(`Updated: '${selectedValue}'`)
        // console.log(select.innerHTML)
    });
}
    

function generateSchema() {
    const schema = {
        "$schema": "http://json-schema.org/draft-06/schema#",
        "$ref": "#/definitions/Main",
        "definitions": {}
    };

    document.querySelectorAll(".type-section").forEach(section => {
        const typeName = section.getAttribute("data-type");
        const pascalTypeName = toPascalCase(typeName);
        schema.definitions[pascalTypeName] = {
            "type": "object",
            "additionalProperties": false,
            "properties": {},
            "required": [],
            "title": pascalTypeName
        };

        section.querySelectorAll(".property-section").forEach(propSection => {
            const propName = propSection.querySelector(".property-name").value;
            const propType = propSection.querySelector(".property-type").value;
            const isArray = propSection.querySelector(".property-array").checked;

            if (propName) {
                let property;
                if (jsonTypes.includes(propType)) {
                    property = isArray ? { "type": "array", "items": { "type": propType } } : { "type": propType };
                } else {
                    const pascalPropType = toPascalCase(propType);
                    property = isArray ? { "type": "array", "items": { "$ref": `#/definitions/${pascalPropType}` } } : { "$ref": `#/definitions/${pascalPropType}` };
                }
                schema.definitions[pascalTypeName].properties[propName] = property;
                schema.definitions[pascalTypeName].required.push(propName);
            }
        });
    });

    const schemaTxt = JSON.stringify(schema, null, 2);
    document.getElementById("json-schema").textContent = schemaTxt;
    return schemaTxt;
}

function initializePageWithSchema(schemaString) {
    const schema = JSON.parse(schemaString);
    const definitions = schema.definitions;

    customTypes = [];

    document.getElementById("type-definitions").innerHTML = '';

    for (const typeName in definitions) {
        if (definitions.hasOwnProperty(typeName)) {
            const originalTypeName = typeName;
            const originalTypeNameCamel = originalTypeName.replace(/([A-Z])/g, ' $1').trim();

            addType(originalTypeNameCamel);

            const typeSection = document.querySelector(`.type-section[data-type="${originalTypeNameCamel}"]`);
            const properties = definitions[typeName].properties;
            console.log(properties);
            for (const propName in properties) {
                if (properties.hasOwnProperty(propName)) {
                    const prop = properties[propName];
                    const isArray = prop.type === "array";
                    const propType = isArray ? (prop.items.$ref ? prop.items.$ref.replace('#/definitions/', '') : prop.items.type) : (prop.$ref ? prop.$ref.replace('#/definitions/', '') : prop.type);
                    // console.log(propType)
                    addProperty(typeSection.querySelector("button"), propName, propType, isArray);
                }
            }
        }
    }

    // TODO: commenting this kinda fixes the initial problem, 
    // but pre existing checkboxes do not have all types
    // updateTypeOptions();
}
