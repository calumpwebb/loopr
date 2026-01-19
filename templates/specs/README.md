# Specifications

This directory contains feature specifications for the Ralph Loop to implement.

## Purpose

Specifications serve as the single source of truth for what features should exist and how they should behave. The planning and build loops read these specs to understand what to implement.

## Structure

Each spec should be a markdown file describing a feature or component:

```
specs/
├── README.md (this file)
├── example-spec.md
├── authentication.md
├── api-endpoints.md
└── ...
```

## Writing Specs

Good specifications are:

1. **Clear** - Unambiguous about what needs to be implemented
2. **Complete** - Include all requirements, edge cases, and constraints
3. **Testable** - Describe how to verify the feature works
4. **Independent** - Can be implemented without reading other specs (or reference them explicitly)

## Template

See `example-spec.md` for a template you can copy and adapt.

## Spec Lifecycle

1. **Author** - Write spec describing desired feature
2. **Plan** - Run `loopr plan` to analyze specs and create implementation plan
3. **Build** - Run `loopr build` to implement according to specs
4. **Update** - If inconsistencies found, agent may update specs (requires Opus 4.5)

## Best Practices

- Keep specs focused on WHAT, not HOW (let the agent decide implementation)
- Include examples of expected behavior
- Document why a feature exists (context helps agents make better decisions)
- Update specs when requirements change
- Remove or archive specs for deprecated features
