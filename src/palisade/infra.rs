use std::collections::HashSet;
use std::option::Option;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

impl Ident {
    pub fn as_str(&self) -> String {
        format!("{}::{}", &self.package, &self.name)
    }

    pub fn new(package: &str, name: &str) -> Self {
        Ident {
            package: package.to_string(),
            name: name.to_string(),
        }
    }
}

pub fn base_ident(name: &str) -> Ident {
    Ident {
        package: "primary".to_string(),
        name: name.to_string(),
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Application {
    pub alpha: Option<Box<Expression>>,
    pub app: Ident,
    pub omega: Box<Expression>,
}

impl Application {
    fn as_str(&self) -> String {
        format!(
            "({}{} {})",
            match &self.alpha {
                Some(a) => format!("{} ", a.as_str()),
                None => "".to_string(),
            },
            self.app.as_str(),
            self.omega.as_str(),
        )
    }

    pub fn expr(&self) -> Expression {
        Expression::Application(self.clone())
    }
}

#[derive(Debug, Clone, PartialEq, Hash, Eq)]
pub enum AtomicType {
    Bool,
    Char,
    Int,
    Real,
    Void,
}

impl Morpheme {
    fn as_str(&self) -> String {
        match self {
            Morpheme::Bool(x) => x.to_string(),
            Morpheme::Char(x) => format!("'{}'", (*x as char).to_string()),
            Morpheme::Int(x) => x.to_string(),
            Morpheme::Real(x) => x.to_string(),
            Morpheme::Void => "Void".to_string(),
        }
    }

    pub fn expr(&self) -> Expression {
        Expression::Morpheme(*self)
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Morpheme {
    Bool(bool),
    Char(u8),
    Int(i64),
    Real(f64),
    Void,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Vector {
    pub body: Vec<Expression>,
}

impl Vector {
    fn as_str(&self) -> String {
        format!(
            "{}",
            self.body
                .iter()
                .map(|x| x.as_str())
                .collect::<Vec<_>>()
                .join(" ")
        )
    }

    pub fn expr(&self) -> Expression {
        Expression::Vector(self.clone())
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum Expression {
    Application(Application),
    Vector(Vector),
    Morpheme(Morpheme),
}

impl Expression {
    pub fn as_str(&self) -> String {
        match self {
            Expression::Application(x) => x.as_str(),
            Expression::Vector(x) => x.as_str(),
            Expression::Morpheme(x) => x.as_str(),
        }
    }

    pub fn expr(&self) -> Expression {
        match self {
            Expression::Application(x) => x.expr(),
            Expression::Vector(x) => x.expr(),
            Expression::Morpheme(x) => x.expr(),
        }
    }
}

pub struct Function {
    pub ident: Ident,
    pub alpha: Option<TypeGroup>,
    pub omega: TypeGroup,
    pub sigma: TypeGroup,
    pub body: Option<Vec<Box<Expression>>>,
}

impl Function {
    pub fn as_str(&self) -> String {
        format!(
            "{}{} {} -> {}\n\t{}",
            match &self.alpha {
                Some(x) => format!("{} ", x.as_str()),
                None => "".to_string(),
            },
            self.ident.as_str(),
            self.omega.as_str(),
            self.sigma.as_str(),
            match &self.body {
                Some(x) => x
                    .iter()
                    .map(|x| x.as_str())
                    .collect::<Vec<_>>()
                    .join("\n\t"),
                None => "".to_string(),
            },
        )
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct TypeGroup {
    pub gamma: HashSet<Type>,
    void: bool,
    universal: bool,
}

impl TypeGroup {
    pub fn universal(&self) -> bool {
        self.universal
    }

    pub fn void(&self) -> bool {
        self.void
    }

    pub fn as_str(&self) -> String {
        format!(
            "{{{}}}",
            self.gamma
                .iter()
                .map(|x| x.as_str())
                .collect::<Vec<_>>()
                .join(" ")
        )
    }

    pub fn of(ts: &Type) -> TypeGroup {
        let mut gamma = HashSet::new();
        gamma.insert(ts.clone());
        TypeGroup {
            gamma: gamma,
            void: false,
            universal: false,
        }
    }

    pub fn of_void() -> TypeGroup {
        TypeGroup {
            gamma: HashSet::new(),
            void: true,
            universal: false,
        }
    }

    pub fn of_universal() -> TypeGroup {
        TypeGroup {
            gamma: HashSet::new(),
            void: false,
            universal: true,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Hash, Eq)]
pub enum Type {
    Unknown(Ident),
    Atomic(AtomicType),
    Vector(Box<Type>),
}

impl Type {
    pub fn as_str(&self) -> String {
        match self {
            Type::Unknown(x) => x.as_str(),
            Type::Atomic(x) => x.as_str(),
            Type::Vector(x) => format!("[{}]", x.as_str()),
        }
    }
}

impl AtomicType {
    pub fn as_str(&self) -> String {
        match self {
            AtomicType::Bool => "Bool",
            AtomicType::Char => "Char",
            AtomicType::Int => "Int",
            AtomicType::Real => "Real",
            AtomicType::Void => "Void",
        }
        .to_string()
    }
}