import clsx from "clsx";
import Heading from "@theme/Heading";
import styles from "./styles.module.css";

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<"svg">>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: "Easy to Implement",
    Svg: require("@site/static/img/easy-to-implement.svg").default,
    description: (
      <>
        Our message broker is designed for straightforward integration, allowing
        you to set up and start using it with minimal effort.
      </>
    ),
  },
  {
    title: "Safe Event Processing",
    Svg: require("@site/static/img/safe-event-processing.svg").default,
    description: (
      <>
        Our message broker ensures reliable event insertion and processing by
        leveraging PostgreSQL's ACID transactions, guaranteeing data integrity
        and consistency.
      </>
    ),
  },
  {
    title: "Flexible and Scalable",
    Svg: require("@site/static/img/flexible-and-scalable.svg").default,
    description: (
      <>
        Built with scalability in mind, our message broker adapts to your
        growing needs while maintaining high performance and flexibility in
        design.
      </>
    ),
  },
];

function Feature({ title, Svg, description }: FeatureItem) {
  return (
    <div className={clsx("col col--4")}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
