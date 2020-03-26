# Text panels

> The text panel lets you make information and description panels etc. for your
> dashboards.
>
> — https://grafana.com/docs/grafana/latest/features/panels/text/

```yaml
rows:
  - name: "Text panels row"
    panels:
      - text:
          title: Some awesome text?
          span: 6
          height: 400px
          markdown: "Markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)\n${percentile}"

      - text:
          title: Some awesome html?
          span: 3
          height: 200px
          html: "Some <b>awesome</b> html?"
```

## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
