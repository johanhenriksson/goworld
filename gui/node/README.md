

Render process


For each node src, dst:
    - Compare type. 
      If wrong type:
        - Disregard src and continue with dst
      Otherwise:
        - Compare and update props src props
    - Render/expand node, i.e. use the props to render the child node
        - If component:
            Components are subtrees
            - Render returns a new node
            - This node is the only child!
            - What is the widget? Component widget!!
        - If basic element:
            Basic elements are leaf nodes in the VDOM
            - Hydrate
            - Set widget
    - Update node
        - Set widges
        - Set children
    - Recursively go through children