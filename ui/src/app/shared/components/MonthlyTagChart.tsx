import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import {
  Chart,
  ChartAxis,
  ChartLine,
  ChartContainer,
  ChartThemeColor,
  ChartLegend,
  ChartVoronoiContainer,
  ChartGroup,
} from '@patternfly/react-charts/victory';
import { IMonthlyTagReport } from '@app/shared/models/tag';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface MonthlyTagChartProps {
  data: IMonthlyTagReport[];
  loading?: boolean;
  title?: string;
  tagNames?: string[];
  startDate?: string;
  endDate?: string;
}

const MonthlyTagChart: React.FC<MonthlyTagChartProps> = ({
  data,
  loading = false,
  title = 'Monthly Tag Amounts',
  tagNames = [],
  startDate,
  endDate,
}) => {
  const { theme } = useTheme();
  // Validate date range spans multiple months
  const dateRangeError = React.useMemo(() => {
    if (!startDate || !endDate) return null;

    const start = new Date(startDate);
    const end = new Date(endDate);

    // Check if dates span multiple months
    const startMonthYear = `${start.getFullYear()}-${start.getMonth()}`;
    const endMonthYear = `${end.getFullYear()}-${end.getMonth()}`;

    if (startMonthYear === endMonthYear) {
      return 'Date range must span multiple months to display monthly trends';
    }

    return null;
  }, [startDate, endDate]);

  const chartData = React.useMemo(() => {
    if (!data || data.length === 0 || tagNames.length === 0) return {};
    if (dateRangeError) return {};

    // Filter data by date range if provided
    let filteredData = data;
    if (startDate && endDate) {
      const start = new Date(startDate);
      const end = new Date(endDate);

      filteredData = data.filter((item) => {
        const itemDate = new Date(item.year, item.month - 1); // month is 1-based in data
        return itemDate >= start && itemDate <= end;
      });
    }

    // Group data by tag
    const tagGroups: { [tagName: string]: Array<{ name: string; x: string; y: number }> } = {};

    tagNames.forEach((tagName) => {
      const tagData = filteredData
        .filter((item) => item.tag === tagName)
        .sort((a, b) => {
          if (a.year !== b.year) return a.year - b.year;
          return a.month - b.month;
        })
        .map((item) => ({
          x: `${item.year}-${String(item.month).padStart(2, '0')}`,
          y: item.amount,
          name: tagName,
        }));

      if (tagData.length > 0) {
        tagGroups[tagName] = tagData;
      }
    });

    return tagGroups;
  }, [data, tagNames, startDate, endDate, dateRangeError]);

  // Calculate the maximum Y value from all chart data for dynamic domain
  const maxYValue = React.useMemo(() => {
    const allValues: number[] = [];
    Object.values(chartData).forEach((tagData) => {
      tagData.forEach((point) => {
        allValues.push(point.y);
      });
    });

    if (allValues.length === 0) return 1000; // Default fallback

    const max = Math.max(...allValues);
    // Add 10% padding to the max value for better visualization
    return Math.ceil(max * 1.1);
  }, [chartData]);

  // Extract all unique x values from chartData for x-axis
  const xAxisValues = React.useMemo(() => {
    const allXValues = new Set<string>();
    Object.values(chartData).forEach((tagData) => {
      tagData.forEach((point) => {
        allXValues.add(point.x);
      });
    });

    // Sort the x values chronologically
    return Array.from(allXValues).sort();
  }, [chartData]);

  // Create legend data from selected tags
  const legendData = React.useMemo(() => {
    return tagNames.map((tagName) => ({
      name: tagName,
    }));
  }, [tagNames]);

  if (loading) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <Spinner size="lg" />
          </div>
        </CardBody>
      </Card>
    );
  }

  // Show error message if date range is invalid
  if (dateRangeError) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div style={{ textAlign: 'center', color: 'var(--pf-v6-global--danger-color--100)' }}>
              <strong>Error:</strong> {dateRangeError}
            </div>
          </div>
        </CardBody>
      </Card>
    );
  }

  if (Object.keys(chartData).length === 0) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div>No data available</div>
          </div>
        </CardBody>
      </Card>
    );
  }

  return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
                                       <div style={{ 
            height: '600px',
            width: '100%'
          }}>
        <Chart
          ariaDesc="Monthly tag amounts chart"
          ariaTitle={title}
          containerComponent={
            <ChartVoronoiContainer 
              labels={({ datum }) => `${datum.name}: ${datum.y.toLocaleString('fr-FR', {
                style: 'currency',
                currency: 'EUR',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })}`} 
              constrainToVisibleArea 
            />
          }
          legendData={legendData}
          legendOrientation="vertical"
          legendPosition="right"
          maxDomain={{ y: maxYValue }}
          minDomain={{ y: 0 }}
          name="monthlyTagChart"
          padding={{
            bottom: 80,
            left: 80,
            right: 120,
            top: 50,
          }}
        >
         <ChartAxis
           tickValues={xAxisValues}
           tickFormat={(value) => {
             // Format month-year string for display (e.g., "2024-01" -> "Jan 2024")
             const [year, month] = value.split('-');
             const monthNames = [
               'Jan',
               'Feb',
               'Mar',
               'Apr',
               'May',
               'Jun',
               'Jul',
               'Aug',
               'Sep',
               'Oct',
               'Nov',
               'Dec',
             ];
             return `${monthNames[parseInt(month) - 1]} ${year}`;
           }}
         />
                   <ChartAxis dependentAxis showGrid tickFormat={(value) =>
            value.toLocaleString('fr-FR', {
              style: 'currency',
              currency: 'EUR',
              minimumFractionDigits: 0,
              maximumFractionDigits: 0,
            })
          } />
         <ChartGroup>
           {Object.entries(chartData).map(([tagName, tagData]) => (
             <ChartLine key={tagName} data={tagData} />
           ))}
                  </ChartGroup>
       </Chart>
     </div>
     {tagNames.length > 0 && (
       <div
         style={{
           textAlign: 'center',
           marginTop: '1rem',
           fontSize: '0.875rem',
           color: 'var(--pf-v6-global--Color--200)',
         }}
       >
         Showing amounts for: <strong>{tagNames.join(', ')}</strong>
       </div>
     )}
     </CardBody>
     </Card>
  );
};

export { MonthlyTagChart };
