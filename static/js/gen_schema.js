let idCounter = 0;

function addChild(parentId) {
    const parent = parentId === 'main' ? document.getElementById('children-container') : document.getElementById(`children-${parentId}`);
    
    idCounter++;
    const newChildId = `child-${idCounter}`;
    
    const childDiv = document.createElement('div');
    childDiv.classList.add('child');
    childDiv.id = newChildId;
    
    const nameInput = document.createElement('input');
    nameInput.type = 'text';
    nameInput.placeholder = 'Child Name';
    nameInput.id = `name-${newChildId}`;
    
    const addButton = document.createElement('button');
    addButton.textContent = 'Add Child';
    addButton.onclick = () => addChild(newChildId);
    
    const childrenContainer = document.createElement('div');
    childrenContainer.id = `children-${newChildId}`;
    
    childDiv.appendChild(nameInput);
    childDiv.appendChild(addButton);
    childDiv.appendChild(childrenContainer);
    
    parent.appendChild(childDiv);
}

function generateJSON() {
    const mainName = document.getElementById('main-name').value;
    const children = buildChildren(document.getElementById('children-container'));
    
    const result = {
        name: mainName,
        children: children
    };
    
    document.getElementById('output').textContent = JSON.stringify(result, null, 2);
}

function buildChildren(container) {
    const children = [];
    const childDivs = container.children;
    
    for (const childDiv of childDivs) {
        const childName = childDiv.querySelector('input').value;
        const grandChildrenContainer = childDiv.querySelector('div');
        const grandChildren = buildChildren(grandChildrenContainer);
        
        children.push({
            name: childName,
            children: grandChildren
        });
    }
    
    return children;
}

